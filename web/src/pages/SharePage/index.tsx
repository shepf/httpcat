import {
  CheckCircleFilled,
  ClockCircleOutlined,
  CloudDownloadOutlined,
  ExclamationCircleFilled,
  FileTextOutlined,
  LockFilled,
  PictureOutlined,
  UnlockFilled,
  UserOutlined,
} from '@ant-design/icons';
import { Button, Input, message, Spin, Typography } from 'antd';
import { useEffect, useState } from 'react';
import { history } from 'umi';
import { getShareInfo, verifyShareCode } from '@/services/ant-design-pro/api';

const { Text } = Typography;

// ===== 内联样式 =====

const pageStyle: React.CSSProperties = {
  minHeight: '100vh',
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  padding: '20px',
  background: '#0f0f1a',
  position: 'relative',
  overflow: 'hidden',
};

const bgOrb = (top: string, left: string, size: number, color: string, delay: string): React.CSSProperties => ({
  position: 'absolute',
  top,
  left,
  width: size,
  height: size,
  borderRadius: '50%',
  background: color,
  filter: 'blur(80px)',
  opacity: 0.5,
  animation: `float 8s ease-in-out ${delay} infinite alternate`,
});

const cardStyle: React.CSSProperties = {
  width: '100%',
  maxWidth: 440,
  background: 'rgba(255, 255, 255, 0.05)',
  backdropFilter: 'blur(24px)',
  WebkitBackdropFilter: 'blur(24px)',
  borderRadius: 24,
  border: '1px solid rgba(255, 255, 255, 0.1)',
  padding: '40px 36px',
  position: 'relative',
  zIndex: 1,
  boxShadow: '0 24px 48px rgba(0, 0, 0, 0.4), inset 0 1px 0 rgba(255,255,255,0.06)',
};

const SharePage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [shareInfo, setShareInfo] = useState<API.ShareInfoResult | null>(null);
  const [extractCode, setExtractCode] = useState('');
  const [verified, setVerified] = useState(false);
  const [verifying, setVerifying] = useState(false);
  const [downloadHover, setDownloadHover] = useState(false);
  const [downloading, setDownloading] = useState(false);
  const [downloadError, setDownloadError] = useState<string | null>(null);
  const [downloadSuccess, setDownloadSuccess] = useState(false);

  const pathParts = history.location.pathname.split('/');
  const shareCode = pathParts[pathParts.length - 1] || '';

  useEffect(() => {
    if (!shareCode) {
      setLoading(false);
      return;
    }
    fetchShareInfo();
  }, [shareCode]);

  const fetchShareInfo = async () => {
    setLoading(true);
    try {
      const res = await getShareInfo(shareCode);
      setShareInfo(res);
      if (res.valid && res.share && !res.share.hasExtractCode) {
        setVerified(true);
      }
    } catch {
      setShareInfo({ valid: false, reason: '获取分享信息失败' });
    } finally {
      setLoading(false);
    }
  };

  const handleVerify = async () => {
    if (!extractCode.trim()) {
      message.warning('请输入提取码');
      return;
    }
    setVerifying(true);
    try {
      const res = await verifyShareCode(shareCode, extractCode.trim());
      if (res.valid) {
        setVerified(true);
        message.success('验证成功');
      } else {
        message.error(res.reason || '提取码错误');
      }
    } catch {
      message.error('验证失败');
    } finally {
      setVerifying(false);
    }
  };

  const handleDownload = async () => {
    setDownloadError(null);
    setDownloading(true);
    try {
      const codeParam = shareInfo?.share?.hasExtractCode ? `?code=${encodeURIComponent(extractCode)}` : '';
      const url = `/s/${shareCode}/download${codeParam}`;
      const resp = await fetch(url);

      // 后端返回错误（403、404、500 等），尝试解析 JSON 错误信息
      if (!resp.ok) {
        let errMsg = '下载失败';
        try {
          const errData = await resp.json();
          errMsg = errData.error || errMsg;
        } catch {
          // 非 JSON 响应，使用默认错误信息
        }
        setDownloadError(errMsg);
        return;
      }

      // 成功：触发浏览器下载
      const blob = await resp.blob();
      const disposition = resp.headers.get('Content-Disposition');
      let fileName = share.fileName;
      if (disposition) {
        const match = disposition.match(/filename\*?=(?:UTF-8''|"?)([^";]+)/i);
        if (match) fileName = decodeURIComponent(match[1].replace(/"/g, ''));
      }
      const blobUrl = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = blobUrl;
      a.download = fileName;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(blobUrl);

      // 下载成功：标记成功 & 本地更新下载计数（不重新请求后端，避免 valid=false 导致页面变成"不可用"）
      setDownloadSuccess(true);
      if (shareInfo?.share) {
        setShareInfo({
          ...shareInfo,
          share: {
            ...shareInfo.share,
            curDownloads: (shareInfo.share.curDownloads || 0) + 1,
          },
        });
      }
    } catch {
      setDownloadError('网络错误，下载失败');
    } finally {
      setDownloading(false);
    }
  };

  const getFileExt = (name: string) => name.split('.').pop()?.toUpperCase() || '';
  const isImage = (type?: string) => type === 'image';

  // 全局 keyframes
  const globalStyles = `
    @keyframes float {
      0% { transform: translate(0, 0) scale(1); }
      100% { transform: translate(30px, -30px) scale(1.1); }
    }
    @keyframes fadeInUp {
      from { opacity: 0; transform: translateY(20px); }
      to { opacity: 1; transform: translateY(0); }
    }
    @keyframes pulse {
      0%, 100% { transform: scale(1); }
      50% { transform: scale(1.05); }
    }
    @keyframes shimmer {
      0% { background-position: -200% center; }
      100% { background-position: 200% center; }
    }
  `;

  // ===== Loading =====
  if (loading) {
    return (
      <div style={pageStyle}>
        <style>{globalStyles}</style>
        <div style={bgOrb('10%', '10%', 300, 'rgba(99, 102, 241, 0.4)', '0s')} />
        <div style={bgOrb('60%', '70%', 250, 'rgba(236, 72, 153, 0.35)', '2s')} />
        <Spin size="large" />
      </div>
    );
  }

  // ===== 无效分享 =====
  if (!shareInfo?.valid) {
    return (
      <div style={pageStyle}>
        <style>{globalStyles}</style>
        <div style={bgOrb('10%', '10%', 300, 'rgba(99, 102, 241, 0.4)', '0s')} />
        <div style={bgOrb('60%', '70%', 250, 'rgba(236, 72, 153, 0.35)', '2s')} />
        <div style={{ ...cardStyle, textAlign: 'center', animation: 'fadeInUp 0.6s ease-out' }}>
          <div style={{
            width: 72, height: 72, borderRadius: '50%',
            background: 'rgba(239, 68, 68, 0.15)',
            display: 'flex', alignItems: 'center', justifyContent: 'center',
            margin: '0 auto 20px',
          }}>
            <ExclamationCircleFilled style={{ fontSize: 36, color: '#ef4444' }} />
          </div>
          <div style={{ fontSize: 22, fontWeight: 600, color: '#fff', marginBottom: 8 }}>
            分享不可用
          </div>
          <div style={{ fontSize: 14, color: 'rgba(255,255,255,0.5)', marginBottom: 32 }}>
            {shareInfo?.reason || '该分享链接已失效或不存在'}
          </div>
          <Button
            type="primary"
            onClick={() => history.push('/')}
            style={{
              height: 44, borderRadius: 12, fontSize: 15, fontWeight: 500,
              background: 'rgba(255,255,255,0.1)', border: '1px solid rgba(255,255,255,0.15)',
              color: '#fff',
            }}
          >
            返回首页
          </Button>
        </div>
      </div>
    );
  }

  const share = shareInfo.share!;
  const ext = getFileExt(share.fileName);

  return (
    <div style={pageStyle}>
      <style>{globalStyles}</style>

      {/* 背景装饰球 */}
      <div style={bgOrb('5%', '5%', 320, 'rgba(99, 102, 241, 0.45)', '0s')} />
      <div style={bgOrb('65%', '75%', 280, 'rgba(236, 72, 153, 0.35)', '2s')} />
      <div style={bgOrb('30%', '60%', 200, 'rgba(14, 165, 233, 0.3)', '4s')} />

      <div style={{ ...cardStyle, animation: 'fadeInUp 0.6s ease-out' }}>

        {/* ===== 顶部品牌 ===== */}
        <div style={{ textAlign: 'center', marginBottom: 32 }}>
          <div style={{
            display: 'inline-flex', alignItems: 'center', gap: 8,
            padding: '6px 16px', borderRadius: 20,
            background: 'rgba(255,255,255,0.06)', border: '1px solid rgba(255,255,255,0.08)',
          }}>
            <span style={{ fontSize: 16 }}>🐱</span>
            <span style={{ fontSize: 13, fontWeight: 500, color: 'rgba(255,255,255,0.7)', letterSpacing: 1 }}>
              HttpCat
            </span>
          </div>
        </div>

        {/* ===== 文件图标 + 名称 ===== */}
        <div style={{ textAlign: 'center', marginBottom: 28 }}>
          <div style={{
            width: 80, height: 80, borderRadius: 20,
            background: isImage(share.fileType)
              ? 'linear-gradient(135deg, rgba(236,72,153,0.2), rgba(99,102,241,0.2))'
              : 'linear-gradient(135deg, rgba(99,102,241,0.2), rgba(14,165,233,0.2))',
            border: `1px solid ${isImage(share.fileType) ? 'rgba(236,72,153,0.3)' : 'rgba(99,102,241,0.3)'}`,
            display: 'flex', alignItems: 'center', justifyContent: 'center',
            margin: '0 auto 16px',
          }}>
            {isImage(share.fileType)
              ? <PictureOutlined style={{ fontSize: 32, color: '#ec4899' }} />
              : <FileTextOutlined style={{ fontSize: 32, color: '#818cf8' }} />
            }
          </div>

          <div style={{
            fontSize: 18, fontWeight: 600, color: '#fff',
            marginBottom: 6, wordBreak: 'break-all', lineHeight: 1.4,
          }}>
            {share.fileName}
          </div>

          {ext && (
            <span style={{
              display: 'inline-block', padding: '2px 10px', borderRadius: 6,
              background: 'rgba(99,102,241,0.15)', color: '#a5b4fc',
              fontSize: 11, fontWeight: 600, letterSpacing: 1,
            }}>
              {ext}
            </span>
          )}
        </div>

        {/* ===== 文件信息卡片 ===== */}
        <div style={{
          background: 'rgba(255,255,255,0.04)',
          borderRadius: 14, padding: '16px 20px', marginBottom: 28,
          border: '1px solid rgba(255,255,255,0.06)',
        }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', flexWrap: 'wrap', gap: 12 }}>
            {/* 分享者 */}
            <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
              <div style={{
                width: 28, height: 28, borderRadius: 8,
                background: 'rgba(99,102,241,0.15)',
                display: 'flex', alignItems: 'center', justifyContent: 'center',
              }}>
                <UserOutlined style={{ fontSize: 13, color: '#818cf8' }} />
              </div>
              <Text style={{ color: 'rgba(255,255,255,0.6)', fontSize: 13 }}>
                {share.createdBy}
              </Text>
            </div>

            {/* 过期时间 */}
            {share.expireAt && (
              <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <div style={{
                  width: 28, height: 28, borderRadius: 8,
                  background: 'rgba(251,191,36,0.12)',
                  display: 'flex', alignItems: 'center', justifyContent: 'center',
                }}>
                  <ClockCircleOutlined style={{ fontSize: 13, color: '#fbbf24' }} />
                </div>
                <Text style={{ color: 'rgba(255,255,255,0.6)', fontSize: 13 }}>
                  {new Date(share.expireAt).toLocaleDateString('zh-CN', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })}
                </Text>
              </div>
            )}

            {/* 下载次数 */}
            <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
              <div style={{
                width: 28, height: 28, borderRadius: 8,
                background: 'rgba(34,197,94,0.12)',
                display: 'flex', alignItems: 'center', justifyContent: 'center',
              }}>
                <CloudDownloadOutlined style={{ fontSize: 13, color: '#22c55e' }} />
              </div>
              <Text style={{ color: 'rgba(255,255,255,0.6)', fontSize: 13 }}>
                {share.maxDownloads && share.maxDownloads > 0
                  ? `${share.curDownloads}/${share.maxDownloads} 次`
                  : `${share.curDownloads} 次`
                }
              </Text>
            </div>
          </div>
        </div>

        {/* ===== 提取码输入 ===== */}
        {share.hasExtractCode && !verified && (
          <div style={{ marginBottom: 24 }}>
            <div style={{
              display: 'flex', alignItems: 'center', gap: 8,
              marginBottom: 16, color: 'rgba(255,255,255,0.5)', fontSize: 13,
            }}>
              <LockFilled style={{ color: '#fbbf24' }} />
              <span>该文件需要提取码才能下载</span>
            </div>
            <div style={{ display: 'flex', gap: 10 }}>
              <Input
                placeholder="请输入提取码"
                value={extractCode}
                onChange={(e) => setExtractCode(e.target.value)}
                onPressEnter={handleVerify}
                maxLength={10}
                style={{
                  flex: 1, height: 48, borderRadius: 12,
                  background: 'rgba(255,255,255,0.06)',
                  border: '1px solid rgba(255,255,255,0.1)',
                  color: '#fff', fontSize: 18, fontWeight: 600,
                  letterSpacing: 8, textAlign: 'center',
                }}
              />
              <Button
                type="primary"
                loading={verifying}
                onClick={handleVerify}
                style={{
                  height: 48, width: 80, borderRadius: 12,
                  background: 'linear-gradient(135deg, #6366f1, #8b5cf6)',
                  border: 'none', fontWeight: 600, fontSize: 15,
                }}
              >
                验证
              </Button>
            </div>
          </div>
        )}

        {/* ===== 验证成功 / 直接下载 ===== */}
        {verified && (
          <div>
            {share.hasExtractCode && (
              <div style={{
                display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 8,
                marginBottom: 20, padding: '10px 0',
                color: '#34d399', fontSize: 14,
              }}>
                <UnlockFilled />
                <span>提取码验证通过</span>
                <CheckCircleFilled />
              </div>
            )}

            {/* ===== 下载成功提示 ===== */}
            {downloadSuccess && (
              <div style={{
                display: 'flex', alignItems: 'center', gap: 10,
                padding: '12px 16px', marginBottom: 16, borderRadius: 12,
                background: 'rgba(34, 197, 94, 0.1)',
                border: '1px solid rgba(34, 197, 94, 0.25)',
                animation: 'fadeInUp 0.3s ease-out',
              }}>
                <CheckCircleFilled style={{ fontSize: 18, color: '#22c55e', flexShrink: 0 }} />
                <span style={{ color: '#86efac', fontSize: 14, lineHeight: 1.5 }}>
                  文件下载成功！
                  {share.maxDownloads > 0 && share.curDownloads >= share.maxDownloads && (
                    <span style={{ color: 'rgba(255,255,255,0.4)', marginLeft: 6 }}>
                      （该分享下载次数已用完）
                    </span>
                  )}
                </span>
              </div>
            )}

            {/* ===== 下载错误提示 ===== */}
            {downloadError && (
              <div style={{
                display: 'flex', alignItems: 'center', gap: 10,
                padding: '12px 16px', marginBottom: 16, borderRadius: 12,
                background: 'rgba(239, 68, 68, 0.1)',
                border: '1px solid rgba(239, 68, 68, 0.25)',
                animation: 'fadeInUp 0.3s ease-out',
              }}>
                <ExclamationCircleFilled style={{ fontSize: 18, color: '#ef4444', flexShrink: 0 }} />
                <span style={{ color: '#fca5a5', fontSize: 14, lineHeight: 1.5 }}>
                  {downloadError}
                </span>
              </div>
            )}

            {(() => {
              const exhausted = share.maxDownloads > 0 && share.curDownloads >= share.maxDownloads;
              return (
                <Button
                  type="primary"
                  size="large"
                  icon={<CloudDownloadOutlined />}
                  onClick={handleDownload}
                  loading={downloading}
                  disabled={exhausted}
                  onMouseEnter={() => setDownloadHover(true)}
                  onMouseLeave={() => setDownloadHover(false)}
                  block
                  style={{
                    height: 52, borderRadius: 14, fontSize: 16, fontWeight: 600,
                    background: exhausted
                      ? 'rgba(255,255,255,0.08)'
                      : downloadHover
                        ? 'linear-gradient(135deg, #818cf8, #a78bfa)'
                        : 'linear-gradient(135deg, #6366f1, #8b5cf6)',
                    border: 'none',
                    boxShadow: exhausted ? 'none' : '0 8px 24px rgba(99, 102, 241, 0.35)',
                    transition: 'all 0.3s ease',
                    transform: downloadHover && !exhausted ? 'translateY(-1px)' : 'none',
                    color: exhausted ? 'rgba(255,255,255,0.3)' : '#fff',
                  }}
                >
                  {downloading ? '下载中...' : exhausted ? '下载次数已用完' : '下载文件'}
                </Button>
              );
            })()}
          </div>
        )}

        {/* ===== 底部 ===== */}
        <div style={{
          textAlign: 'center', marginTop: 28, paddingTop: 20,
          borderTop: '1px solid rgba(255,255,255,0.05)',
          color: 'rgba(255,255,255,0.25)', fontSize: 12,
        }}>
          Powered by HttpCat · 安全文件分享
        </div>
      </div>
    </div>
  );
};

export default SharePage;
