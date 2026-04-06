import ChangePasswordModal from '@/pages/user/Login/ChangePasswordModal';
import { changePasswd, outLogin } from '@/services/ant-design-pro/api';
import { LockOutlined, LogoutOutlined } from '@ant-design/icons';
import { Avatar, Form, Menu, message, Spin } from 'antd';
import type { ItemType } from 'antd/lib/menu/hooks/useItems';
import { stringify } from 'querystring';
import type { MenuInfo } from 'rc-menu/lib/interface';
import React, { useCallback, useState } from 'react';
import { history, useModel } from 'umi';
import HeaderDropdown from '../HeaderDropdown';
import styles from './index.less';

export type GlobalHeaderRightProps = {
  menu?: boolean;
};

const loginOut = async () => {
  await outLogin();
  localStorage.removeItem('token');

  const { query = {}, search, pathname } = history.location;
  const { redirect } = query;
  if (window.location.pathname !== '/user/login' && !redirect) {
    history.replace({
      pathname: '/user/login',
      search: stringify({ redirect: pathname + search }),
    });
  }
};

const AvatarDropdown: React.FC<GlobalHeaderRightProps> = () => {
  const { initialState, setInitialState } = useModel('@@initialState');
  const [modalVisible, setModalVisible] = useState(false);
  const [form] = Form.useForm();

  const handleCancelModal = () => {
    form.resetFields();
    setModalVisible(false);
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();
      if (values.newPassword !== values.confirmPassword) {
        message.error('两次输入的密码不一致');
        return;
      }

      const { oldPassword, newPassword } = values;
      const response = await changePasswd({ oldPassword, newPassword });

      if (response.errorCode === 0) {
        message.success('密码修改成功');
        form.resetFields();
        localStorage.removeItem('token');
        setModalVisible(false);
        setInitialState((s) => ({ ...s, currentUser: undefined }));

        const { query = {}, search, pathname } = history.location;
        const { redirect } = query;
        if (window.location.pathname !== '/user/login' && !redirect) {
          history.replace({
            pathname: '/user/login',
            search: stringify({ redirect: pathname + search }),
          });
        }
      } else if (response.errorCode === 2) {
        message.error('旧密码错误');
      } else {
        message.error(response.msg || '密码修改失败');
      }
    } catch (error) {
      message.error('密码修改失败');
    }
  };

  const onMenuClick = useCallback(
    (event: MenuInfo) => {
      const { key } = event;
      if (key === 'logout') {
        setInitialState((s) => ({ ...s, currentUser: undefined }));
        loginOut();
        return;
      }
      if (key === 'changePasswd') {
        setModalVisible(true);
        return;
      }
      history.push(`/account/${key}`);
    },
    [setInitialState],
  );

  const loading = (
    <span className={`${styles.action} ${styles.account}`}>
      <Spin size="small" style={{ marginLeft: 8, marginRight: 8 }} />
    </span>
  );

  if (!initialState) return loading;
  const { currentUser } = initialState;
  if (!currentUser || !currentUser.name) return loading;

  const menuItems: ItemType[] = [
    {
      key: 'changePasswd',
      icon: <LockOutlined />,
      label: '修改密码',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
    },
  ];

  const menuHeaderDropdown = (
    <Menu className={styles.menu} selectedKeys={[]} onClick={onMenuClick} items={menuItems} />
  );

  return (
    <div>
      <HeaderDropdown overlay={menuHeaderDropdown}>
        <span className={`${styles.action} ${styles.account}`}>
          <Avatar size="small" className={styles.avatar} src={currentUser.avatar} alt="avatar" />
          <span className={`${styles.name} anticon`}>{currentUser.name}</span>
        </span>
      </HeaderDropdown>
      <ChangePasswordModal visible={modalVisible} onCancel={handleCancelModal} onOk={handleModalOk} form={form} />
    </div>
  );
};

export default AvatarDropdown;
