import { getConf, getDownloadStatistics, getUploadAvailableSpace, getUploadStatistics, getVersion } from '@/services/ant-design-pro/api';
import { Bar } from '@ant-design/charts';
import { ProCard, ProDescriptions, Statistic, StatisticCard } from '@ant-design/pro-components';
import { useEffect, useState } from 'react';
import { Pie } from '@ant-design/plots';
import { CheckCircleTwoTone, CloseCircleTwoTone } from '@ant-design/icons';
import { Space } from 'antd';



const imgStyle = {
  display: 'block',
  width: 42,
  height: 42,
};

// 组件的主要结构是使用 ProCard 组件进行布局管理，并在其中嵌套了各种统计卡片和图表组件。
export default () => {

    const [responsive, setResponsive] = useState(false);
    const [collapsed, setCollapsed] = useState(false);


    //获取配置信息
    const [confData, setConfData] = useState({
      downloadDir: '',
      fileUploadEnable: true,
      uploadDir: '',
      webDir: '',
    });
    useEffect(() => {
      // 在组件挂载时获取信息
      const fetchConf = async () => {
        try {
          const response  = await getConf();
          const data = response.data;
          // 处理版本号为空的情况，设置为空字符串
          const confDataToUpdate = {
            downloadDir: data.downloadDir || '',
            fileUploadEnable: data.fileUploadEnable || true,
            uploadDir: data.uploadDir || '',
            webDir: data.webDir || '',
          };
          setConfData(confDataToUpdate);
        } catch (error) {
          console.error('Failed to fetch conf:', error);
        }
      };
  
      fetchConf();
    }, []);


    // 获取版本信息
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


    // 获取上传统计信息
    const [uploadStatisticsData, setUploadStatisticsData] = useState({
      monthPercentage: '',
      monthUploadCount: 0,
      lastMonthUploadCount: 0,
      todayPercentage: '',
      todayUploadCount: 0,
      yesterdayUploadCount: 0,
      totalUploadCount: 0,
    });
    useEffect(() => {
      // 在组件挂载时获取信息
      const fetchUploadStatistics = async () => {
        try {
          const response  = await getUploadStatistics();
          const data = response.data;
          // 处理为空的情况，设置为空字符串
          const dataToUpdate = {
            monthPercentage: data.monthPercentage || '',
            monthUploadCount: data.monthUploadCount || 0,
            lastMonthUploadCount: data.lastMonthUploadCount || 0,
            todayPercentage: data.todayPercentage || '',
            todayUploadCount: data.todayUploadCount || 0,
            yesterdayUploadCount: data.yesterdayUploadCount || 0,
            totalUploadCount: data.totalUploadCount || 0,
          };
          setUploadStatisticsData(dataToUpdate);
        } catch (error) {
          console.error('Failed to fetch UploadStatistics:', error);
        }
      };
  
      fetchUploadStatistics();
    }, []);
  



    // 获取上传统计信息
    const [downloadStatisticsData, setDownloadStatisticsData] = useState({
      monthPercentage: '',
      monthDownloadCount: 0,
      lastMonthDownloadCount: 0,
      todayPercentage: '',
      todayDownloadCount: 0,
      yesterdayDownloadCount: 0,
      totalDownloadCount: 0,
    });
    useEffect(() => {
      // 在组件挂载时获取信息
      const fetchDownloadStatistics = async () => {
        try {
          const response  = await getDownloadStatistics();
          const data = response.data;
          // 处理为空的情况，设置为空字符串
          const downloadDataToUpdate = {
            monthPercentage: data.monthPercentage || '',
            monthDownloadCount: data.monthDownloadCount || 0,
            lastMonthDownloadCount: data.lastMonthDownloadCount || 0,
            todayPercentage: data.todayPercentage || '',
            todayDownloadCount: data.todayDownloadCount || 0,
            yesterdayDownloadCount: data.yesterdayDownloadCount || 0,
            totalDownloadCount: data.totalDownloadCount || 0,
          };
          setDownloadStatisticsData(downloadDataToUpdate);
        } catch (error) {
          console.error('Failed to fetch DownloadStatistics:', error);
        }
      };
  
      fetchDownloadStatistics();
    }, []);


    // 获取磁盘信息
    const [usedSpace, setUsedSpace] = useState(0);
    const [freeSpace, setFreeSpace] = useState(0);
    const totalSize = parseFloat(((usedSpace + freeSpace) / (1024 * 1024 * 1024)).toFixed(2)); // 总大小，单位：GB
    // toFixed() 方法返回的是一个字符串类型的值，而不是一个数字
    const usedSize = parseFloat((usedSpace/(1024 * 1024 * 1024)).toFixed(2)); // 使用多少，单位：GB
    const availSize = parseFloat((freeSpace/(1024 * 1024 * 1024)).toFixed(2)); // 剩余多少空间，单位：GB
    const usagePercentage = (usedSize / totalSize * 100).toFixed(1) + '%';
    useEffect(() => {
      // 在组件挂载时获取版本信息
      const fetchUploadAvailableSpace = async () => {
        try {
          const data = await getUploadAvailableSpace();
          setUsedSpace(data.usedSpace || 0); 
          setFreeSpace(data.freeSpace || 0); 


        } catch (error) {
          console.error('Failed to fetch Upload AvailableSpace:', error);
        }
      };
  
      fetchUploadAvailableSpace();
    }, []);


    // 我们使用普通饼图实现
    const DiskInfoPie = () => {

        const data = [
        {
          type: '已用',
          value: usedSize,
          // value: 10,
        },
        {
          type: '剩余',
          value: availSize,
          // value: 20,
        },
      ];
      const config = {
        data,
        angleField: 'value',
        colorField: 'type',
        radius: 0.8,
        label: {
          text: (d: { type: any; value: any; }) => `${d.type}\n ${d.value}G`,
          style: {
            fontWeight: 'bold',
          },
          // position: 'spider',
        },
        legend: {
          color: {
            title: false,
            position: 'right',
            rowPadding: 5,
          },
        },
      };
      return <Pie {...config} />;
    };
    

    // // 图例 有问题，只能显示 已用，TODO，暂时不使用环形饼图
    // const DiskInfoPie = () => {

    //   const data = [
    //     {
    //       type: '已用',
    //       value: usedSize,
    //     },
    //     {
    //       type: '空闲',
    //       value: availSize,
    //     },
    //   ];

    //   const config = {
    //     data: data,
    //     angleField: 'value',
    //     colorField: 'type', // 使用类型字段进行颜色映射
    //     paddingRight: 80,
    //     innerRadius: 0.6,
    //     label: {
    //       text: (d: { type: any; value: any; }) => `${d.type}\n ${d.value}G`,
    //       position: 'spider',
    //       style: {
    //         fontWeight: 'bold',
    //       },
    //     },
    //     // label: {
    //     //   text: 'value',
    //     //   style: {
    //     //     fontWeight: 'bold',
    //     //   },
    //     // },
    //     legend: { //图例
    //       color: {
    //         fields: ['type'], // 添加'type'字段
    //         title: false,
    //         position: 'right',
    //         rowPadding: 5,
    //       },
    //     },
    //     annotations: [
    //       {
    //         type: 'text',
    //         style: {
    //           text: `总计:\n ${totalSize} G`,
    //           x: '50%',
    //           y: '50%',
    //           textAlign: 'center',
    //           fontSize: 40,
    //           fontStyle: 'bold',
    //         },
    //       },
    //     ],
      
    //   };
    //   return <Pie {...config} />;
    // };



    // 下载top3
    const FileDownloadBarChart = () => {
      const data = [
        {
          type: '文件afsdfadfafafafa1',
          sales: 145,
        },
        {
          type: '文件2adfafafafafadfaf',
          sales: 61,
        },
        {
          type: '文件3asdfasdfadfaaafaasfafadfadgwergtwerty',
          sales: 52,
        },
      ];
      const config = {
        data,
        xField: 'sales',
        yField: 'type',
        seriesField: 'type',
        legend: false,
        scrollbar: {
          type: 'vertical',
        },
        //指定条形图最小最大宽度
        // minBarWidth: 20,
        maxBarWidth: 20,
        meta: {
          type: {
            alias: '文件名',
          },
          sales: {
            alias: '下载次数',
          },
        },
      };
      return <Bar {...config} />;
    };
        

  return (
    <>
      <ProCard
        title="系统基本配置信息"
        // extra="extra"  //右上角显示内容
        tooltip="系统基本配置暂时不支持界面修改，需要你修改配置文件重启服务生效！"
        // style={{ maxWidth: 300 }}
        headerBordered
      >
        <ProDescriptions
        title="httpcat system info"
        dataSource={{
          version: versionData.version,
          httpcat_uptime: versionData.uptime,

          upload_path: confData.uploadDir,
          download_path: confData.downloadDir,
          web_path: confData.webDir,
          fileUploadEnable: confData.fileUploadEnable,
          
          money2: -12345.33,
          state: 'all',
          switch: true,
          state2: 'open',
        }}
        emptyText={'空'}
        columns={[
          {
            title: 'httcat版本',
            key: 'text',
            dataIndex: 'version',
          },
          // {
          //   title: '上传文件开关状态',
          //   key: 'fileUploadEnable',
          //   dataIndex: 'fileUploadEnable',
          //   valueType: 'select',
          //   valueEnum: {
          //     true: {
          //       text: '开启',
          //       status: 'Error',
          //     },
          //     false: {
          //       text: '关闭',
          //       status: 'Success',
          //     },
          //   },
          // },
          {
            title: '上传文件开关状态',
            key: 'fileUploadEnable',
            dataIndex: 'fileUploadEnable',
            render: (text, record) => (
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

          {
            title: '操作',
            valueType: 'option',
            render: () => [
              <a target="_blank" rel="noopener noreferrer" key="link">
                {/* 使用帮助 */}
              </a>,
            ],
          },
        ]}
      >
        {/* <ProDescriptions.Item
          dataIndex="percent"
          label="上传路径空间剩余"
          valueType="percent"
        >
          30
        </ProDescriptions.Item>
        <ProDescriptions.Item label="超链接">
          <a href="alipay.com">超链接</a>
        </ProDescriptions.Item> */}
        </ProDescriptions>
      </ProCard>

      <ProCard
        title="数据概览"
        extra="" //用于右上角显示
        split={responsive ? 'horizontal' : 'vertical'}
        headerBordered
        bordered
        // collapsible
        ghost
        // gutter={8} 
      >
        <ProCard split="horizontal">
          <ProCard split="horizontal">
            <ProCard split="vertical">
              <StatisticCard
                statistic={{
                  title: '今日上传文件个数',
                  value: uploadStatisticsData.todayUploadCount,
                  description: (
                    <>
                      <p  style={{marginBottom: 0}}>昨日上传: {uploadStatisticsData.yesterdayUploadCount}</p>
                      <Statistic
                        title="较昨日"
                        value= {uploadStatisticsData.todayPercentage}
                        trend= {uploadStatisticsData.todayPercentage.startsWith('-') ? 'down' : 'up'}
                      />
                    </>
                  ),
                }}
              />
              <StatisticCard
                statistic={{
                  title: '今日下载文件个数',
                  value: downloadStatisticsData.todayDownloadCount,
                  description: (
                    <>
                       <p  style={{marginBottom: 0}}>昨日下载: {downloadStatisticsData.yesterdayDownloadCount}</p>
                       <Statistic
                        title="较昨日"
                        value= {downloadStatisticsData.todayPercentage}
                        trend= {downloadStatisticsData.todayPercentage.startsWith('-') ? 'down' : 'up'}
                        />
                    </>
                  ),
                }}
              />
            </ProCard>
            <ProCard split="vertical">
            <StatisticCard
                statistic={{
                  title: '本月累计上传文件个数',
                  value: uploadStatisticsData.monthUploadCount,
                  description: (
                    <>
                      <span> 上月上传: {uploadStatisticsData.lastMonthUploadCount} </span>
                      <Statistic title="月同比" value={uploadStatisticsData.monthPercentage} trend={uploadStatisticsData.monthPercentage.startsWith('-') ? 'down' : 'up'} />
                    </>
                  ),
                }}
              />
              <StatisticCard
                statistic={{
                  title: '本月累计下载文件个数',
                  value: downloadStatisticsData.monthDownloadCount,
                  description: (
                    <>
                      <span> 上月下载: {downloadStatisticsData.lastMonthDownloadCount} </span>
                      <Statistic title="月同比" value={downloadStatisticsData.monthPercentage} trend={downloadStatisticsData.monthPercentage.startsWith('-') ? 'down' : 'up'} />
                    </>
                  ),
                }}
              />
            </ProCard>
            <ProCard split="vertical">
              <StatisticCard
                statistic={{
                  title: '总计上传文件个数',
                  value: uploadStatisticsData.totalUploadCount,
                  suffix: '个',
                }}
              />
              <StatisticCard
                statistic={{
                  title: '总计下载文件个数',
                  value: downloadStatisticsData.totalDownloadCount,
                  suffix: '个',
                }}
              />
            </ProCard>

            {/* <ProCard split="vertical">
            <StatisticCard
              statistic={{
                title: '下载文件top3',
              }}
              chart={<FileDownloadBarChart />}
              />
            </ProCard> */}

          </ProCard>

        </ProCard>
        <StatisticCard
              statistic={{
                title: '上传目录空间',
                value: `已使用${usedSize}G，剩余${availSize}G`,
                description: <Statistic title="已使用占比" value={usagePercentage} />,
              }}
              chart={<DiskInfoPie />}
            
        />

      </ProCard>

    </>


  );
};