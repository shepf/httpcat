import React, { useState } from 'react';
import { Button, Card, Modal, Input, Form, message } from 'antd';

type FormValues = {
  oldPassword: string;
  newPassword: string;
  confirmPassword: string;
};

type ChangePasswordModalProps = {
  visible: boolean;
  onCancel: () => void;
  onOk: () => void;
  form: any; // 添加 form 属性
};

const ChangePasswordModal: React.FC<ChangePasswordModalProps> = ({ visible, onCancel, onOk, form }) => {


  return (
    <Modal
      title="修改密码"
      open={visible}
      onCancel={onCancel}
      onOk={onOk}
      maskClosable={false}  // 不希望用户点击背景时默认取消弹框，您可以将 maskClosable 属性设置为 false
    >
      <Form form={form}>
        <Form.Item
          name="oldPassword"
          label="旧密码"
          rules={[{ required: true, message: '请输入旧密码' }]}
        >
          <Input.Password />
        </Form.Item>
        <Form.Item
          name="newPassword"
          label="新密码"
          rules={[{ required: true, message: '请输入新密码' }]}
        >
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