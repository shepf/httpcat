#!/bin/bash
# 搜索图片并上传到 httpcat
# 从 Unsplash 已验证可用的技术类图片中随机下载并上传到 httpcat 图片管理
# 用法: ./search-and-upload-image.sh "关键词" [数量]
# 示例: ./search-and-upload-image.sh "技术文章封面" 2

QUERY="$1"; COUNT=${2:-2}
SD="$(cd "$(dirname "$0")" && pwd)"
TD="${SD}/temp_images"
[ -z "$QUERY" ] && echo "用法: $0 <关键词> [数量]" >&2 && exit 1
mkdir -p "$TD" 2>/dev/null || true

# 已验证可用的 unsplash 技术类图片 ID 列表
IDS=(
  '1518770660439-4636190af475'
  '1488590528505-98d2b5aba04b'
  '1461749280684-dccba630e2f6'
  '1504639725590-34d0984388bd'
  '1526374965328-7f61d4dc18c5'
  '1550751827-4bd374c3f58b'
  '1555949963-ff9fe0c870eb'
  '1498050108023-c5249f4df085'
  '1517694712202-14dd9538aa97'
  '1519389950473-47ba0277781c'
  '1531297484001-80022131f5a1'
)

echo "[1/2] 下载图片..."
# 随机选取
SHUF_IDS=($(printf '%s\n' "${IDS[@]}" | shuf | head -n $COUNT))
DL=0; FILES=""
for pid in "${SHUF_IDS[@]}"; do
  DL=$((DL+1))
  FN="${QUERY// /_}_${DL}.jpg"; FP="${TD}/${FN}"
  URL="https://images.unsplash.com/photo-${pid}?w=1200&fit=crop&q=80"
  echo "  下载: $URL"
  HC=$(curl -sL --max-time 20 -o "$FP" -w "%{http_code}" "$URL" 2>/dev/null)
  FS=$(stat -c%s "$FP" 2>/dev/null || echo 0)
  if [ "$HC" = "200" ] && [ "$FS" -gt 5000 ]; then
    echo "  OK: $FN (${FS}B)"; FILES="${FILES}${FP}\n"
  else echo "  FAIL: HTTP$HC ${FS}B" >&2; rm -f "$FP"; fi
done
[ -z "$FILES" ] && echo "没有下载到图片" >&2 && exit 1

echo "[2/2] 上传到httpcat..."
OK=0
while IFS= read -r fp; do
  [ -z "$fp" ] && continue
  echo "  上传: $fp"
  R=$(bash "${SD}/httpcat-api.sh" UPLOAD_IMAGE "$fp" 2>&1)
  echo "  $R"
  echo "$R" | grep -qi "upload successful" && OK=$((OK+1))
done < <(echo -e "$FILES")
echo "完成! 上传${OK}张到httpcat图片管理"
rm -f ${TD}/* 2>/dev/null
