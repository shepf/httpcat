import { getVersion } from '@/services/ant-design-pro/api';
import { GithubOutlined } from '@ant-design/icons';
import { DefaultFooter } from '@ant-design/pro-components';
import { useEffect, useState } from 'react';
import { useIntl } from 'umi';

const Footer: React.FC = () => {
  const intl = useIntl();
  const defaultMessage = intl.formatMessage({
    id: 'app.copyright.produced',
    defaultMessage: '',
  });


  const softwareName = 'httpcat';
  // useState 的参数可以是结构体或对象，而不仅限于字符串
  const [versionData, setVersionData] = useState({
    build: '',
    ci: '',
    commit: '',
    uptime: '',
    version: '',
  });



  useEffect(() => {
    // 在组件挂载时获取版本信息
    const fetchVersion = async () => {
      try {
        const response  = await getVersion();
        const data = response.data;
        // 处理版本号为空的情况，设置为空字符串
        const versionDataToUpdate = {
          build: data.build || '',
          ci: data.ci || '',
          commit: data.commit || '',
          uptime: data.uptime || '',
          version: data.version || '',
        };
        setVersionData(versionDataToUpdate);
      } catch (error) {
        console.error('Failed to fetch version:', error);
      }
    };

    fetchVersion();
  }, []);

  
  const currentYear = new Date().getFullYear();

  return (
    <DefaultFooter
      copyright={`${currentYear} ${softwareName} ${versionData.version} ${defaultMessage}`}
      links={[
        // {
        //   key: 'Ant Design Pro',
        //   title: 'Ant Design Pro',
        //   href: 'https://pro.ant.design',
        //   blankTarget: true,
        // },
        {
          key: 'github',
          title: <GithubOutlined />,
          href: 'https://github.com/shepf/httpcat-release/releases',
          blankTarget: true,
        },
        // {
        //   key: 'Ant Design',
        //   title: 'Ant Design',
        //   href: 'https://ant.design',
        //   blankTarget: true,
        // },
      ]}
    />
  );
};

export default Footer;
