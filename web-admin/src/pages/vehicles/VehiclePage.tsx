import { useState, useEffect, useCallback } from 'react';
import {
  Table, Button, Card, Tag, Space, Typography, Modal, Form, Input, Select, message, Popconfirm,
} from 'antd';
import { PlusOutlined, EditOutlined } from '@ant-design/icons';
import apiClient from '../../api/client';

interface VehicleRecord {
  id: number;
  plate: string;
  type: number;
  typeName: string;
  status: number;
  statusName: string;
  fuelLevel: number;
}

interface TypeOption {
  value: number;
  label: string;
}

const STATUS_MAP: Record<number, string> = { 1: '空闲', 2: '装载中', 3: '运输中', 4: '卸载中', 5: '维修中', 6: '离线' };
const STATUS_COLORS: Record<number, string> = { 1: 'green', 2: 'blue', 3: 'processing', 4: 'orange', 5: 'warning', 6: 'default' };

export default function VehiclePage() {
  const [vehicles, setVehicles] = useState<VehicleRecord[]>([]);
  const [typeOptions, setTypeOptions] = useState<TypeOption[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<VehicleRecord | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [form] = Form.useForm();

  const fetchVehicles = useCallback(async () => {
    setLoading(true);
    try {
      const res = await apiClient.get('/api/v1/vehicles');
      const types = typeOptions.length > 0 ? typeOptions : [];
      const getTypeName = (t: number) => types.find((x) => x.value === t)?.label || '未知';
      const list: VehicleRecord[] = (res.data.data || []).map((v: any) => ({
        id: v.id,
        plate: v.plate,
        type: v.type,
        typeName: getTypeName(v.type),
        status: v.status,
        statusName: STATUS_MAP[v.status] || '未知',
        fuelLevel: v.fuel_level ?? 100,
      }));
      setVehicles(list);
    } catch {
      message.error('获取车辆列表失败');
    } finally {
      setLoading(false);
    }
  }, [typeOptions]);

  const fetchVehicleTypes = async () => {
    try {
      const res = await apiClient.get('/api/v1/vehicle-types');
      const list: TypeOption[] = (res.data.data || []).map((t: any) => ({
        value: t.id,
        label: t.name,
      }));
      setTypeOptions(list);
    } catch {
      // ignore
    }
  };

  useEffect(() => { fetchVehicles(); fetchVehicleTypes(); }, []);

  const openEdit = (v: VehicleRecord) => {
    setEditing(v);
    form.setFieldsValue({ plate: v.plate, type: v.type, status: v.status });
    setModalOpen(true);
  };

  const openAdd = () => {
    setEditing(null);
    form.resetFields();
    setModalOpen(true);
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      setSubmitting(true);
      let res: any;
      if (editing) {
        res = await apiClient.put(`/api/v1/vehicles/${editing.id}`, {
          plate: values.plate,
          type: values.type,
          status: values.status,
          mineId: 1,
        });
      } else {
        res = await apiClient.post('/api/v1/vehicles', {
          plate: values.plate,
          type: values.type,
          mineId: 1,
        });
      }
      if (res.data.code === 0 || res.data.message === 'success') {
        message.success(editing ? '更新成功' : '添加成功');
        setModalOpen(false);
        form.resetFields();
        setEditing(null);
        fetchVehicles();
      } else {
        message.error(res.data.message || '操作失败');
      }
    } catch (err: any) {
      if (err?.response?.data?.message) {
        message.error(err.response.data.message);
      } else if (err?.errorFields) {
        // form validation error, antd shows inline
      } else {
        message.error('操作失败');
      }
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = async (id: number) => {
    try {
      const res = await apiClient.delete(`/api/v1/vehicles/${id}`);
      if (res.data.code === 0 || res.data.message === 'success') {
        message.success('删除成功');
        fetchVehicles();
      }
    } catch {
      message.error('删除失败');
    }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 100 },
    { title: '车牌/编号', dataIndex: 'plate', key: 'plate' },
    { title: '类型', dataIndex: 'typeName', key: 'type' },
    {
      title: '状态', dataIndex: 'statusName', key: 'status',
      render: (_: string, record: VehicleRecord) => (
        <Tag color={STATUS_COLORS[record.status]}>{record.statusName}</Tag>
      ),
    },
    {
      title: '油量', dataIndex: 'fuelLevel', key: 'fuelLevel',
      render: (v: number) => `${v}%`,
    },
    {
      title: '操作', key: 'action',
      render: (_: any, record: VehicleRecord) => (
        <Space>
          <a onClick={() => openEdit(record)}><EditOutlined /> 编辑</a>
          <Popconfirm title="确认删除该车辆？" onConfirm={() => handleDelete(record.id)}>
            <a style={{ color: 'red' }}>删除</a>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Typography.Title level={4}>车辆管理</Typography.Title>
      <Card>
        <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'flex-end' }}>
          <Button type="primary" icon={<PlusOutlined />} onClick={openAdd}>
            添加车辆
          </Button>
        </div>
        <Table columns={columns} dataSource={vehicles} rowKey="id" loading={loading} />
      </Card>

      <Modal
        title={editing ? '编辑车辆' : '添加车辆'}
        open={modalOpen}
        onCancel={() => { setModalOpen(false); setEditing(null); }}
        onOk={handleSubmit}
        confirmLoading={submitting}
      >
        <Form form={form} layout="vertical">
          <Form.Item label="车牌/编号" name="plate" rules={[{ required: true, message: '请输入车牌' }]}>
            <Input placeholder="请输入车牌" />
          </Form.Item>
          <Form.Item label="车辆类型" name="type" rules={[{ required: true, message: '请选择车辆类型' }]}>
            <Select placeholder="请选择" options={typeOptions.length > 0 ? typeOptions : [
              { value: 1, label: '矿用卡车' },
              { value: 2, label: '挖掘机' },
              { value: 3, label: '装载机' },
              { value: 4, label: '推土机' },
            ]} />
          </Form.Item>
          {editing && (
            <Form.Item label="状态" name="status">
              <Select options={[
                { value: 1, label: '空闲' },
                { value: 2, label: '装载中' },
                { value: 3, label: '运输中' },
                { value: 4, label: '卸载中' },
                { value: 5, label: '维修中' },
                { value: 6, label: '离线' },
              ]} />
            </Form.Item>
          )}
        </Form>
      </Modal>
    </div>
  );
}
