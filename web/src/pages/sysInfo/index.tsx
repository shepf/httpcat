import {
  getConf,
  getDownloadStatistics,
  getUploadAvailableSpace,
  getUploadStatistics,
  getVersion,
} from '@/services/ant-design-pro/api';
import { ProCard, ProDescriptions, Statistic, StatisticCard } from '@ant-design/pro-components';
import { useEffect, useState } from 'react';
import { Pie } from '@ant-design/plots';
import { CheckCircleTwoTone, CloseCircleTwoTone } from '@ant-design/icons';
import { Space, Spin } from 'antd';

export default () => {
  const [loading, setLoading] = useState(true);
  const [confData, setConfData] = useState<API.HttpcatConf>({});
  const [versionData, setVersionData] = useState<API.Version>({});
  const [uploadStats, setUploadStats] = useState<API.UploadStatistics>({});
  const [downloadStats, setDownloadStats] = useState<API.DownloadStatistics>({});
  const [usedSpace, setUsedSpace] = useState(0);
  const [freeSpace, setFreeSpace] = useState(0);

  // 合并所有数据加载到单个 useEffect
  useEffect(() => {
    const fetchAllData = async () => {
      try {
        const [confRes, versionRes, uploadRes, downloadRes, spaceRes] = await Promise.all([
          getConf(),
          getVersion(),
          getUploadStatistics(),
          getDownloadStatistics(),
          getUploadAvailableSpace(),
        ]);

        setConfData(confRes.data || {});
        setVersionData(versionRes.data || {});
        setUploadStats(uploadRes.data || {});
        setDownloadStats(downloadRes.data || {});
        setUsedSpace(spaceRes.usedSpace || 0);
        setFreeSpace(spaceRes.freeSpace || 0);
      } catch (error) {
        console.error('获取系统信息失败:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchAllData();
  }, []);

  const totalSize = parseFloat(((usedSpace + freeSpace) / (1024 * 1024 * 1024)).toFixed(2));
  const usedSize = parseFloat((usedSpace / (1024 * 1024 * 1024)).toFixed(2));
  const availSize = parseFloat((freeSpace / (1024 * 1024 * 1024)).toFixed(2));
  const usagePercentage = totalSize > 0 ? (usedSize / totalSize * 100).toFixed(1) + '%' : '0%';

  const DiskInfoPie = () => {
    const data = [
      { type: '已用', value: usedSize },
      { type: '剩余', value: availSize },
    ];
    const config = {
      data,
      angleField: 'value',
      colorField: 'type',
      radius: 0.8,
      label: {
        text: (d: { type: string; value: number }) => `${d.type}\n ${d.value}G`,
        style: { fontWeight: 'bold' },
      },
      legend: {
        color: {
          title: false,
          position: 'right' as const,
          rowPadding: 5,
        },
      },
    };
    return <Pie {...config} />;
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '100px 0' }}>
        <Spin size="large" tip="加载系统信息..." />
      </div>
    );
  }

  const todayUploadTrend = (uploadStats.todayPercentage || '0%').startsWith('-') ? 'down' : 'up';
  const todayDownloadTrend = (downloadStats.todayPercentage || '0%').startsWith('-') ? 'down' : 'up';
  const monthUploadTrend = (uploadStats.monthPercentage || '0%').startsWith('-') ? 'down' : 'up';
  const monthDownloadTrend = (downloadStats.monthPercentage || '0%').startsWith('-') ? 'down' : 'up';

  return (
    <>
      <ProCard
        title="系统基本配置信息"
        tooltip="系统基本配置暂时不支持界面修改，需要修改配置文件重启服务生效"
        headerBordered
      >
        <ProDescriptions
          title="httpcat system info"
          dataSource={{
            version: versionData.version,
            httpcat_uptime: versionData.uptime,
            work_dir: confData.workDir,
            upload_path: confData.absUploadDir || confData.uploadDir,
            download_path: confData.absDownloadDir || confData.downloadDir,
            web_path: confData.absWebDir || confData.webDir,
            fileUploadEnable: confData.fileUploadEnable,
          }}
          emptyText="空"
          columns={[
            {
              title: 'httpcat版本',
              key: 'version',
              dataIndex: 'version',
            },
            {
              title: '上传文件开关状态',
              key: 'fileUploadEnable',
              dataIndex: 'fileUploadEnable',
              render: (_text, record) => (
                <Space>
                  {record.fileUploadEnable ? (
                    <CheckCircleTwoTone twoToneColor="#52c41a" />
                  ) : (
                    <CloseCircleTwoTone twoToneColor="#eb2f96" />
                  )}
                  {record.fileUploadEnable ? '开启' : '关闭'}
                </Space>
              ),
            },
            {
              title: 'httpcat持续运行时间',
              key: 'httpcat_uptime',
              dataIndex: 'httpcat_uptime',
            },
            {
              title: '项目工作目录',
              key: 'work_dir',
              dataIndex: 'work_dir',
              copyable: true,
            },
            {
              title: '上传文件路径',
              key: 'upload_path',
              dataIndex: 'upload_path',
              copyable: true,
            },
            {
              title: '下载文件路径',
              key: 'download_path',
              dataIndex: 'download_path',
              copyable: true,
            },
            {
              title: 'web前端路径',
              key: 'web_path',
              dataIndex: 'web_path',
              copyable: true,
            },
          ]}
        />
      </ProCard>

      <ProCard
        title="数据概览"
        split="vertical"
        headerBordered
        bordered
        style={{ marginTop: 16 }}
      >
        <ProCard split="horizontal" colSpan="50%">
          <ProCard split="vertical">
            <StatisticCard
              statistic={{
                title: '今日上传文件个数',
                value: uploadStats.todayUploadCount || 0,
                description: (
                  <>
                    <p style={{ marginBottom: 0 }}>昨日上传: {uploadStats.yesterdayUploadCount || 0}</p>
                    <Statistic title="较昨日" value={uploadStats.todayPercentage || '0%'} trend={todayUploadTrend} />
                  </>
                ),
              }}
            />
            <StatisticCard
              statistic={{
                title: '今日下载文件个数',
                value: downloadStats.todayDownloadCount || 0,
                description: (
                  <>
                    <p style={{ marginBottom: 0 }}>昨日下载: {downloadStats.yesterdayDownloadCount || 0}</p>
                    <Statistic title="较昨日" value={downloadStats.todayPercentage || '0%'} trend={todayDownloadTrend} />
                  </>
                ),
              }}
            />
          </ProCard>
          <ProCard split="vertical">
            <StatisticCard
              statistic={{
                title: '本月累计上传文件个数',
                value: uploadStats.monthUploadCount || 0,
                description: (
                  <>
                    <span>上月上传: {uploadStats.lastMonthUploadCount || 0}</span>
                    <Statistic title="月同比" value={uploadStats.monthPercentage || '0%'} trend={monthUploadTrend} />
                  </>
                ),
              }}
            />
            <StatisticCard
              statistic={{
                title: '本月累计下载文件个数',
                value: downloadStats.monthDownloadCount || 0,
                description: (
                  <>
                    <span>上月下载: {downloadStats.lastMonthDownloadCount || 0}</span>
                    <Statistic title="月同比" value={downloadStats.monthPercentage || '0%'} trend={monthDownloadTrend} />
                  </>
                ),
              }}
            />
          </ProCard>
          <ProCard split="vertical">
            <StatisticCard
              statistic={{
                title: '总计上传文件个数',
                value: uploadStats.totalUploadCount || 0,
                suffix: '个',
              }}
            />
            <StatisticCard
              statistic={{
                title: '总计下载文件个数',
                value: downloadStats.totalDownloadCount || 0,
                suffix: '个',
              }}
            />
          </ProCard>
        </ProCard>
        <StatisticCard
          colSpan="50%"
          statistic={{
            title: '上传目录空间',
            value: `已使用${usedSize}G，剩余${availSize}G`,
            description: <Statistic title="已使用占比" value={usagePercentage} />,
          }}
          chart={<div style={{ height: 280 }}><DiskInfoPie /></div>}
        />
      </ProCard>
    </>
  );
};
