#!/bin/bash
# httpcat AK/SK REST API 客户端
# 用法: ./httpcat-api.sh <METHOD> <PATH> [QUERY] [BODY]
# 示例:
#   ./httpcat-api.sh GET /api/v1/file/listFiles "dir=/"
#   ./httpcat-api.sh POST /api/v1/user/createUploadToken '' '{"appkey":"httpcat","appsecret":"httpcat_app_secret"}'
#   ./httpcat-api.sh UPLOAD /path/to/local/file.txt

set -e

# ── httpcat 配置（根据实际情况修改）──
HOST="http://172.17.0.1:8888"        # API 调用地址（容器内访问 httpcat 用，Docker 网桥地址）
PUBLIC_URL="http://你的公网IP:8888"   # 外部访问地址（返回给用户的下载链接用）
AK="your-access-key"                  # 对应 svr.yml 中 aksk 的 key
SK="your-secret-key"                  # 对应 svr.yml 中 aksk 的 value
APP_KEY="httpcat"                     # 对应 svr.yml 中的 app_key
APP_SECRET="httpcat_app_secret"       # 对应 svr.yml 中的 app_secret

# ── AK/SK 签名函数 ──
sign_request() {
    local method="$1" path="$2" query="$3" body="$4"
    TIMESTAMP=$(date +%s)
    if [ -n "$body" ]; then
        BODY_HASH=$(printf '%s' "$body" | openssl dgst -sha256 -hex | awk '{print $NF}')
    else
        BODY_HASH=$(printf '' | openssl dgst -sha256 -hex | awk '{print $NF}')
    fi
    SIGN_STR=$(printf '%s\n%s\n%s\n%s\n%s\n%s' "$method" "$path" "$query" "$AK" "$TIMESTAMP" "$BODY_HASH")
    SIGNATURE=$(printf '%s' "$SIGN_STR" | openssl dgst -sha256 -hmac "$SK" -hex | awk '{print $NF}')
}

# ── 上传文件（两步：生成 Token + 上传）──
upload_file() {
    local file_path="$1"
    if [ ! -f "$file_path" ]; then
        echo "错误: 文件不存在: $file_path" >&2
        exit 1
    fi

    # Step 1: 生成 UploadToken
    local token_body="{\"appkey\":\"${APP_KEY}\",\"appsecret\":\"${APP_SECRET}\"}"
    sign_request "POST" "/api/v1/user/createUploadToken" "" "$token_body"
    local token_resp=$(curl -s "${HOST}/api/v1/user/createUploadToken" -X POST \
        -H "Content-Type: application/json" \
        -H "AccessKey: ${AK}" \
        -H "Signature: ${SIGNATURE}" \
        -H "TimeStamp: ${TIMESTAMP}" \
        -d "$token_body")

    local upload_token=$(echo "$token_resp" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('data',''))" 2>/dev/null)
    if [ -z "$upload_token" ]; then
        echo "获取 UploadToken 失败: $token_resp" >&2
        exit 1
    fi
    echo "UploadToken: ${upload_token}" >&2

    # Step 2: 上传文件
    local resp=$(curl -s "${HOST}/api/v1/file/upload" -X POST \
        -H "UploadToken: ${upload_token}" \
        -F "f1=@${file_path}")

    # 输出结果，将内网地址替换为公网地址
    echo "$resp" | sed "s|${HOST}|${PUBLIC_URL}|g"
    local filename=$(basename "$file_path")
    echo "" >&2
    echo "公网下载链接: ${PUBLIC_URL}/api/v1/file/download?filename=${filename}" >&2
}

# ── 上传图片到图片管理（生成缩略图）──
upload_image() {
    local file_path="$1"
    [ ! -f "$file_path" ] && echo "file not found: $file_path" >&2 && exit 1
    local tb="{\"appkey\":\"${APP_KEY}\",\"appsecret\":\"${APP_SECRET}\"}"
    sign_request "POST" "/api/v1/user/createUploadToken" "" "$tb"
    local tr=$(curl -s "${HOST}/api/v1/user/createUploadToken" -X POST \
        -H "Content-Type: application/json" \
        -H "AccessKey: ${AK}" -H "Signature: ${SIGNATURE}" -H "TimeStamp: ${TIMESTAMP}" \
        -d "$tb")
    local ut=$(echo "$tr" | python3 -c "import sys,json; print(json.load(sys.stdin).get('data',''))" 2>/dev/null)
    [ -z "$ut" ] && echo "get token failed: $tr" >&2 && exit 1
    echo "UploadToken: ${ut}" >&2
    local resp=$(curl -s "${HOST}/api/v1/imageManage/upload" -X POST \
        -H "UploadToken: ${ut}" -F "file=@${file_path}")

    # 输出结果，将内网地址替换为公网地址
    echo "$resp" | sed "s|${HOST}|${PUBLIC_URL}|g"
    local filename=$(basename "$file_path")
    echo "" >&2
    echo "公网查看链接: ${PUBLIC_URL}/api/v1/imageManage/download?filename=${filename}" >&2
}

# ── 下载文件 ──
download_file() {
    local filename="$1" output="$2"
    if [ -z "$output" ]; then
        output="./$filename"
    fi
    curl -s -o "$output" "${HOST}/api/v1/file/download?filename=${filename}"
    echo "{\"status\":\"ok\",\"saved_to\":\"${output}\"}"
}

# ── 主逻辑 ──
case "${1^^}" in
    UPLOAD)
        upload_file "$2"
        ;;
    UPLOAD_IMAGE)
        upload_image "$2"
        ;;
    DOWNLOAD)
        download_file "$2" "$3"
        ;;
    GET)
        sign_request "GET" "$2" "$3" ""
        URL="${HOST}${2}"
        [ -n "$3" ] && URL="${URL}?${3}"
        curl -s "$URL" \
            -H "AccessKey: ${AK}" \
            -H "Signature: ${SIGNATURE}" \
            -H "TimeStamp: ${TIMESTAMP}"
        ;;
    POST)
        sign_request "POST" "$2" "$3" "$4"
        URL="${HOST}${2}"
        [ -n "$3" ] && URL="${URL}?${3}"
        curl -s "$URL" -X POST \
            -H "Content-Type: application/json" \
            -H "AccessKey: ${AK}" \
            -H "Signature: ${SIGNATURE}" \
            -H "TimeStamp: ${TIMESTAMP}" \
            -d "$4"
        ;;
    *)
        cat << 'HELP'
httpcat AK/SK API 客户端

用法: ./httpcat-api.sh <命令> [参数...]

命令:
  GET <PATH> [QUERY]         发起 GET 请求
  POST <PATH> [QUERY] [BODY] 发起 POST 请求
  UPLOAD <FILE>              上传普通文件（自动获取 Token）
  UPLOAD_IMAGE <FILE>        上传图片到图片管理（生成缩略图，可在图片管理页面查看）
  DOWNLOAD <FILENAME> [OUT]  下载文件

示例:
  ./httpcat-api.sh GET /api/v1/file/listFiles "dir=/"
  ./httpcat-api.sh GET /api/v1/conf/getVersion
  ./httpcat-api.sh GET /api/v1/file/getFileInfo "filename=test.txt"
  ./httpcat-api.sh GET /api/v1/statistics/getStatistics
  ./httpcat-api.sh GET /api/v1/statistics/getDiskUsage
  ./httpcat-api.sh GET /api/v1/statistics/getUploadHistory "page=1&pageSize=10"
  ./httpcat-api.sh UPLOAD /path/to/file.txt
  ./httpcat-api.sh UPLOAD_IMAGE /path/to/photo.png
  ./httpcat-api.sh DOWNLOAD myfile.txt ./saved.txt
HELP
        exit 1
        ;;
esac
