import { getSysConfig, updateSysConfig, restartServer } from '@/services/ant-design-pro/api';
import {
  ApiOutlined,
  CloudUploadOutlined,
  FolderOpenOutlined,
  PictureOutlined,
  FileTextOutlined,
  GlobalOutlined,
  SaveOutlined,
  ReloadOutlined,
  ExclamationCircleOutlined,
  LockOutlined,
  LoadingOutlined,
  SyncOutlined,
} from '@ant-design/icons';
import {
  Card,
  Form,
  Input,
  InputNumber,
  Switch,
  Button,
  Space,
  message,
  Spin,
  Divider,
  Alert,
  Modal,
  Tag,
  Tooltip,
  Select,
  Row,
  Col,
  Result,
} from 'antd';
import { useEffect, useState, useCallback, useRef } from 'react';

const LOG_LEVELS = [
  { value: -1, label: 'Debug (-1)', color: 'default' },
  { value: 0, label: 'Info (0)', color: 'blue' },
  { value: 1, label: 'Warn (1)', color: 'orange' },
  { value: 2, label: 'Error (2)', color: 'red' },
  { value: 3, label: 'DPanic (3)', color: 'magenta' },
  { value: 4, label: 'Panic (4)', color: 'volcano' },
  { value: 5, label: 'Fatal (5)', color: 'purple' },
];

// 需要重启生效的字段映射
const RESTART_FIELD_MAP: Record<string, string> = {
  uploadDir: '上传子目录',
  downloadDir: '下载子目录',
  httpPort: 'HTTP 端口',
};

export default () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [initialValues, setInitialValues] = useState<API.SysConfig>({});
  // 重启相关状态
  const [restartModalVisible, setRestartModalVisible] = useState(false);
  const [restartPassword, setRestartPassword] = useState('');
  const [restarting, setRestarting] = useState(false);
  const [restartPhase, setRestartPhase] = useState<'idle' | 'saving' | 'restarting' | 'reconnecting' | 'done'>('idle');
  // 保存待提交的变更（在密码验证后使用）
  const pendingChangesRef = useRef<Partial<API.SysConfig> | null>(null);

  const fetchConfig = useCallback(async () => {
    setLoading(true);
    try {
      const res = await getSysConfig();
      if (res.data) {
        const config = res.data;
        setInitialValues(config);
        form.setFieldsValue(config);
      }
    } catch (error) {
      message.error('获取系统配置失败');
    } finally {
      setLoading(false);
    }
  }, [form]);

  useEffect(() => {
    fetchConfig();
  }, [fetchConfig]);

  // 计算变更的字段
  const getChangedFields = (values: API.SysConfig): Partial<API.SysConfig> => {
    const changed: Record<string, any> = {};
    // 只读字段不提交到后端
    const readonlyFields = new Set(['fileBaseDir', 'fullUploadDir', 'fullDownloadDir']);
    const keys = Object.keys(values) as (keyof API.SysConfig)[];
    keys.forEach((key) => {
      if (readonlyFields.has(key)) return;
      if (values[key] !== initialValues[key]) {
        changed[key] = values[key];
      }
    });
    return changed;
  };

  // 等待服务恢复（轮询健康检查）
  const waitForServerReady = (): Promise<boolean> => {
    return new Promise((resolve) => {
      let attempts = 0;
      const maxAttempts = 30; // 最多等 30 秒
      const interval = setInterval(async () => {
        attempts++;
        try {
          const res = await fetch('/api/v1/conf/getVersion', {
            method: 'GET',
            headers: { Authorization: `Bearer ${localStorage.getItem('token') || ''}` },
          });
          if (res.ok) {
            clearInterval(interval);
            resolve(true);
          }
        } catch {
          // 服务还没恢复，继续等
        }
        if (attempts >= maxAttempts) {
          clearInterval(interval);
          resolve(false);
        }
      }, 1000);
    });
  };

  const handleSave = async () => {
    try {
      const values = await form.validateFields();
      const changed = getChangedFields(values);

      if (Object.keys(changed).length === 0) {
        message.info('没有配置变更');
        return;
      }

      const restartChanges = Object.keys(changed).filter((k) => k in RESTART_FIELD_MAP);

      if (restartChanges.length > 0) {
        // 有需要重启的配置 → 弹出密码确认框
        pendingChangesRef.current = changed;
        setRestartPassword('');
        setRestartModalVisible(true);
      } else {
        // 纯热更新
        await doSave(changed, false);
      }
    } catch {
      message.error('请检查表单填写是否正确');
    }
  };

  // 密码确认后执行：保存配置 → 重启服务
  const handleRestartConfirm = async () => {
    if (!restartPassword.trim()) {
      message.warning('请输入管理员密码');
      return;
    }

    const changed = pendingChangesRef.current;
    if (!changed) return;

    setRestarting(true);
    setRestartPhase('saving');

    try {
      // 第一步：保存配置
      const saveRes = await updateSysConfig(changed);
      if (saveRes.errorCode !== 0) {
        message.error(saveRes.msg || '保存配置失败');
        setRestarting(false);
        setRestartPhase('idle');
        return;
      }

      setRestartPhase('restarting');

      // 第二步：调用重启接口（需要密码验证）
      try {
        const restartRes = await restartServer(restartPassword);
        if (restartRes.errorCode !== 0) {
          message.error(restartRes.msg || '重启失败');
          setRestarting(false);
          setRestartPhase('idle');
          return;
        }
      } catch (restartErr: any) {
        // 密码错误（403）
        if (restartErr?.response?.status === 403 || restartErr?.data?.errorCode === 2) {
          message.error('管理员密码错误');
          setRestarting(false);
          setRestartPhase('idle');
          return;
        }
        // 其他错误可能是服务已经开始关闭了，继续等待重连
      }

      setRestartPhase('reconnecting');

      // 第三步：等待服务重新启动
      // 先等 2 秒让旧进程完全关闭
      await new Promise((r) => setTimeout(r, 2000));

      const ready = await waitForServerReady();

      if (ready) {
        setRestartPhase('done');
        setTimeout(() => {
          setRestartModalVisible(false);
          setRestarting(false);
          setRestartPhase('idle');
          pendingChangesRef.current = null;
          message.success('服务已重启，新配置已生效！');
          fetchConfig();
        }, 1500);
      } else {
        setRestartPhase('idle');
        setRestarting(false);
        Modal.warning({
          title: '配置已保存，但服务重连超时',
          content: (
            <div>
              <p>配置已写入文件，服务正在重启中。</p>
              <p>如果页面长时间无法访问，请检查：</p>
              <ul>
                <li>服务是否通过 systemd 或 Docker 管理（支持自动拉起）</li>
                <li>新配置（如端口号）是否正确</li>
              </ul>
              <p>您也可以手动刷新页面重试。</p>
            </div>
          ),
        });
      }
    } catch (error: any) {
      message.error(error?.data?.msg || '操作失败');
      setRestarting(false);
      setRestartPhase('idle');
    }
  };

  // 纯热更新的保存
  const doSave = async (changed: Partial<API.SysConfig>, _hasRestart: boolean) => {
    setSaving(true);
    try {
      const res = await updateSysConfig(changed);
      if (res.errorCode === 0) {
        message.success(res.data?.message || '配置保存成功');
        await fetchConfig();
      } else {
        message.error(res.msg || '保存失败');
      }
    } catch (error: any) {
      message.error(error?.data?.msg || '保存配置失败');
    } finally {
      setSaving(false);
    }
  };

  const handleReset = () => {
    form.setFieldsValue(initialValues);
    message.info('已恢复为当前生效的配置');
  };

  // 重启进度提示
  const getRestartPhaseContent = () => {
    switch (restartPhase) {
      case 'saving':
        return (
          <div style={{ textAlign: 'center', padding: '20px 0' }}>
            <LoadingOutlined style={{ fontSize: 32, color: '#1890ff' }} />
            <p style={{ marginTop: 12, color: '#595959' }}>正在保存配置...</p>
          </div>
        );
      case 'restarting':
        return (
          <div style={{ textAlign: 'center', padding: '20px 0' }}>
            <SyncOutlined spin style={{ fontSize: 32, color: '#faad14' }} />
            <p style={{ marginTop: 12, color: '#595959' }}>正在重启服务...</p>
          </div>
        );
      case 'reconnecting':
        return (
          <div style={{ textAlign: 'center', padding: '20px 0' }}>
            <SyncOutlined spin style={{ fontSize: 32, color: '#1890ff' }} />
            <p style={{ marginTop: 12, color: '#595959' }}>等待服务重新启动...</p>
            <p style={{ color: '#8c8c8c', fontSize: 12 }}>正在检测服务状态，请稍候</p>
          </div>
        );
      case 'done':
        return (
          <Result
            status="success"
            title="重启成功"
            subTitle="新配置已生效"
            style={{ padding: '20px 0' }}
          />
        );
      default:
        return null;
    }
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '100px 0' }}>
        <Spin size="large" tip="加载系统配置..." />
      </div>
    );
  }

  // 计算当前表单中哪些是需要重启的变更（用于弹窗展示）
  const pendingChanges = pendingChangesRef.current || {};
  const restartChanges = Object.keys(pendingChanges).filter((k) => k in RESTART_FIELD_MAP);
  const hotChanges = Object.keys(pendingChanges).filter((k) => !(k in RESTART_FIELD_MAP));

  return (
    <div style={{ maxWidth: 900, margin: '0 auto' }}>
      <Alert
        message="配置说明"
        description="标记为「需重启」的配置项保存后需要重启服务才能生效；其他配置保存后立即生效。文件根目录只能通过配置文件修改，上传/下载目录是文件根目录的子目录。"
        type="info"
        showIcon
        style={{ marginBottom: 16 }}
      />

      <Form
        form={form}
        layout="vertical"
        initialValues={initialValues}
      >
        {/* 存储路径配置 */}
        <Card
          title={
            <Space>
              <FolderOpenOutlined />
              <span>存储路径配置</span>
              <Tag color="warning">需重启</Tag>
            </Space>
          }
          style={{ marginBottom: 16 }}
          size="small"
        >
          {/* 文件根目录（只读） */}
          <Form.Item
            label="文件根目录"
            tooltip="文件上传/下载的根目录，只能通过服务端配置文件修改"
          >
            <Input
              value={initialValues.fileBaseDir || '-'}
              disabled
              addonBefore="只读"
              style={{ maxWidth: 400 }}
            />
          </Form.Item>
          <Row gutter={24}>
            <Col span={12}>
              <Form.Item
                label="上传子目录"
                name="uploadDir"
                rules={[
                  { required: true, message: '请输入上传子目录' },
                  {
                    validator: (_, value) => {
                      if (!value) return Promise.resolve();
                      if (value.startsWith('/')) {
                        return Promise.reject('不能以 / 开头');
                      }
                      if (value.includes('..')) {
                        return Promise.reject('不允许包含 ..');
                      }
                      return Promise.resolve();
                    },
                  },
                ]}
                tooltip="相对于文件根目录的子路径"
                extra={initialValues.fullUploadDir ? `完整路径：${initialValues.fullUploadDir}` : undefined}
              >
                <Input placeholder="website/upload/" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                label="下载子目录"
                name="downloadDir"
                rules={[
                  { required: true, message: '请输入下载子目录' },
                  {
                    validator: (_, value) => {
                      if (!value) return Promise.resolve();
                      if (value.startsWith('/')) {
                        return Promise.reject('不能以 / 开头');
                      }
                      if (value.includes('..')) {
                        return Promise.reject('不允许包含 ..');
                      }
                      return Promise.resolve();
                    },
                  },
                ]}
                tooltip="相对于文件根目录的子路径"
                extra={initialValues.fullDownloadDir ? `完整路径：${initialValues.fullDownloadDir}` : undefined}
              >
                <Input placeholder="download/" />
              </Form.Item>
            </Col>
          </Row>
        </Card>

        {/* HTTP 服务配置 */}
        <Card
          title={
            <Space>
              <GlobalOutlined />
              <span>HTTP 服务配置</span>
              <Tag color="warning">需重启</Tag>
            </Space>
          }
          style={{ marginBottom: 16 }}
          size="small"
        >
          {/* 当前访问端口提示（开发模式下访问端口可能与服务端口不同） */}
          {(() => {
            const accessPort = window.location.port ? parseInt(window.location.port, 10) : (window.location.protocol === 'https:' ? 443 : 80);
            const serverPort = initialValues.httpPort;
            if (serverPort && accessPort !== serverPort) {
              return (
                <Alert
                  type="warning"
                  showIcon
                  style={{ marginBottom: 12 }}
                  message={
                    <span>
                      当前通过端口 <Tag color="blue">{accessPort}</Tag> 访问（开发代理），
                      后端服务实际运行在端口 <Tag color="green">{serverPort}</Tag>。
                      生产环境下两者一致。
                    </span>
                  }
                />
              );
            }
            return null;
          })()}
          <Form.Item
            label="后端服务端口"
            name="httpPort"
            rules={[
              { required: true, message: '请输入端口号' },
              {
                type: 'number',
                min: 1,
                max: 65535,
                message: '端口范围 1-65535',
              },
            ]}
            tooltip="后端 API 服务监听的端口，生产环境下即为访问端口。修改后需要重启服务生效。"
          >
            <InputNumber style={{ width: 200 }} min={1} max={65535} />
          </Form.Item>
        </Card>

        {/* 文件上传策略 */}
        <Card
          title={
            <Space>
              <CloudUploadOutlined />
              <span>文件上传策略</span>
              <Tag color="green">热更新</Tag>
            </Space>
          }
          style={{ marginBottom: 16 }}
          size="small"
        >
          <Row gutter={24}>
            <Col span={12}>
              <Form.Item
                label="文件上传开关"
                name="fileUploadEnable"
                valuePropName="checked"
                tooltip="关闭后所有文件上传接口将被禁用"
              >
                <Switch checkedChildren="开启" unCheckedChildren="关闭" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                label="上传Token校验"
                name="enableUploadToken"
                valuePropName="checked"
                tooltip="开启后外部上传需要提供有效的 UploadToken"
              >
                <Switch checkedChildren="开启" unCheckedChildren="关闭" />
              </Form.Item>
            </Col>
          </Row>
          <Divider dashed style={{ margin: '8px 0 16px' }} />
          <Row gutter={24}>
            <Col span={8}>
              <Form.Item
                label="上传凭证有效期(秒)"
                name="uploadPolicyDeadline"
                tooltip="上传策略的 Token 过期时间，0 表示不限制"
              >
                <InputNumber style={{ width: '100%' }} min={0} placeholder="7200" />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                label={
                  <Tooltip title="0 表示不限制">
                    文件最小值(字节)
                  </Tooltip>
                }
                name="uploadPolicyFSizeMin"
              >
                <InputNumber style={{ width: '100%' }} min={0} placeholder="0" />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item
                label={
                  <Tooltip title="0 表示不限制">
                    文件最大值(字节)
                  </Tooltip>
                }
                name="uploadPolicyFSizeLimit"
              >
                <InputNumber style={{ width: '100%' }} min={0} placeholder="0" />
              </Form.Item>
            </Col>
          </Row>
        </Card>

        {/* 企业微信 Bot 配置 */}
        <Card
          title={
            <Space>
              <ApiOutlined />
              <span>企业微信 Bot 通知</span>
              <Tag color="green">热更新</Tag>
            </Space>
          }
          style={{ marginBottom: 16 }}
          size="small"
        >
          <Form.Item
            label="通知开关"
            name="notifyEnable"
            valuePropName="checked"
            tooltip="开启后文件上传成功会通过企业微信 Webhook 推送通知"
          >
            <Switch checkedChildren="开启" unCheckedChildren="关闭" />
          </Form.Item>
          <Form.Item noStyle shouldUpdate={(prev, cur) => prev.notifyEnable !== cur.notifyEnable}>
            {({ getFieldValue }) =>
              getFieldValue('notifyEnable') ? (
                <Form.Item
                  label="Webhook URL"
                  name="persistentNotifyUrl"
                  rules={[
                    { required: true, message: '开启通知时必须填写 Webhook URL' },
                    {
                      pattern: /^https?:\/\//,
                      message: '必须以 http:// 或 https:// 开头',
                    },
                  ]}
                  tooltip="企业微信群机器人的 Webhook 地址"
                >
                  <Input placeholder="https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx" />
                </Form.Item>
              ) : null
            }
          </Form.Item>
        </Card>

        {/* 缩略图配置 */}
        <Card
          title={
            <Space>
              <PictureOutlined />
              <span>缩略图配置</span>
              <Tag color="green">热更新</Tag>
            </Space>
          }
          style={{ marginBottom: 16 }}
          size="small"
        >
          <Row gutter={24}>
            <Col span={12}>
              <Form.Item
                label="缩略图宽度(px)"
                name="thumbWidth"
                rules={[
                  {
                    type: 'number',
                    min: 50,
                    max: 2000,
                    message: '范围 50-2000',
                  },
                ]}
                tooltip="新上传的图片将按此尺寸生成缩略图，已有缩略图不受影响"
              >
                <InputNumber style={{ width: '100%' }} min={50} max={2000} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                label="缩略图高度(px)"
                name="thumbHeight"
                rules={[
                  {
                    type: 'number',
                    min: 50,
                    max: 2000,
                    message: '范围 50-2000',
                  },
                ]}
              >
                <InputNumber style={{ width: '100%' }} min={50} max={2000} />
              </Form.Item>
            </Col>
          </Row>
        </Card>

        {/* 日志配置 */}
        <Card
          title={
            <Space>
              <FileTextOutlined />
              <span>日志配置</span>
              <Tag color="green">热更新</Tag>
            </Space>
          }
          style={{ marginBottom: 16 }}
          size="small"
        >
          <Form.Item
            label="日志级别"
            name="logLevel"
            tooltip="级别越高输出的日志越少：Debug < Info < Warn < Error < Fatal"
          >
            <Select style={{ width: 300 }} options={LOG_LEVELS.map((l) => ({ value: l.value, label: l.label }))} />
          </Form.Item>
        </Card>

        {/* 操作按钮 */}
        <Card size="small">
          <Space>
            <Button
              type="primary"
              icon={<SaveOutlined />}
              onClick={handleSave}
              loading={saving}
              size="large"
            >
              保存配置
            </Button>
            <Button icon={<ReloadOutlined />} onClick={handleReset} size="large">
              重置
            </Button>
            <Button onClick={fetchConfig} size="large">
              刷新
            </Button>
          </Space>
        </Card>
      </Form>

      {/* 重启确认弹窗（含密码验证） */}
      <Modal
        title={null}
        open={restartModalVisible}
        onCancel={() => {
          if (!restarting) {
            setRestartModalVisible(false);
            setRestartPassword('');
            pendingChangesRef.current = null;
          }
        }}
        footer={restarting ? null : undefined}
        okText="确认保存并重启"
        cancelText="取消"
        onOk={handleRestartConfirm}
        width={520}
        closable={!restarting}
        maskClosable={false}
        confirmLoading={restarting}
      >
        {restarting ? (
          getRestartPhaseContent()
        ) : (
          <div>
            <div style={{ display: 'flex', alignItems: 'center', marginBottom: 16 }}>
              <ExclamationCircleOutlined style={{ color: '#faad14', fontSize: 22, marginRight: 12 }} />
              <span style={{ fontSize: 16, fontWeight: 500 }}>确认保存并重启服务</span>
            </div>

            {/* 热更新配置 */}
            {hotChanges.length > 0 && (
              <div style={{
                background: '#f6ffed',
                border: '1px solid #b7eb8f',
                borderRadius: 6,
                padding: '10px 14px',
                marginBottom: 12,
              }}>
                <div style={{ display: 'flex', alignItems: 'center', marginBottom: 4 }}>
                  <Tag color="green" style={{ margin: 0 }}>保存后立即生效</Tag>
                  <span style={{ marginLeft: 8, color: '#8c8c8c', fontSize: 12 }}>共 {hotChanges.length} 项</span>
                </div>
              </div>
            )}

            {/* 需重启配置 */}
            <div style={{
              background: '#fffbe6',
              border: '1px solid #ffe58f',
              borderRadius: 6,
              padding: '10px 14px',
              marginBottom: 16,
            }}>
              <div style={{ display: 'flex', alignItems: 'center', marginBottom: 6 }}>
                <Tag color="warning" style={{ margin: 0 }}>需重启生效</Tag>
                <span style={{ marginLeft: 8, color: '#8c8c8c', fontSize: 12 }}>共 {restartChanges.length} 项</span>
              </div>
              <ul style={{ margin: '4px 0 0 0', paddingLeft: 20, color: '#595959' }}>
                {restartChanges.map((k) => (
                  <li key={k}>{RESTART_FIELD_MAP[k]}</li>
                ))}
              </ul>
            </div>

            {/* 流程说明 */}
            <Alert
              type="info"
              showIcon
              message="保存后将自动重启服务，新配置立即生效"
              description="流程：保存配置 → 验证密码 → 优雅关闭 → 自动重启（由 systemd / Docker 拉起）"
              style={{ marginBottom: 16 }}
            />

            {/* 密码输入 */}
            <div style={{
              background: '#fafafa',
              border: '1px solid #d9d9d9',
              borderRadius: 6,
              padding: '14px 16px',
            }}>
              <div style={{ marginBottom: 8, fontWeight: 500 }}>
                <LockOutlined style={{ marginRight: 6 }} />
                请输入管理员密码确认操作
              </div>
              <Input.Password
                placeholder="请输入 admin 账户密码"
                value={restartPassword}
                onChange={(e) => setRestartPassword(e.target.value)}
                onPressEnter={handleRestartConfirm}
                size="large"
                autoFocus
              />
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
};
