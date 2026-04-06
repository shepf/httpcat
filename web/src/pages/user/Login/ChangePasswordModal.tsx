import React from 'react';
import { Form, Input, Modal } from 'antd';

type ChangePasswordModalProps = {
  visible: boolean;
  onCancel: () => void;
  onOk: () => void;
  form: any;
};

const ChangePasswordModal: React.FC<ChangePasswordModalProps> = ({ visible, onCancel, onOk, form }) => {
  return (
    <Modal title="修改密码" open={visible} onCancel={onCancel} onOk={onOk} maskClosable={false}>
      <Form form={form} layout="vertical">
        <Form.Item name="oldPassword" label="旧密码" rules={[{ required: true, message: '请输入旧密码' }]}>
          <Input.Password />
        </Form.Item>
        <Form.Item name="newPassword" label="新密码" rules={[{ required: true, message: '请输入新密码' }]}>
          <Input.Password />
        </Form.Item>
        <Form.Item
          name="confirmPassword"
          label="确认密码"
          rules={[{ required: true, message: '请再次输入新密码' }]}
        >
          <Input.Password />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default ChangePasswordModal;
