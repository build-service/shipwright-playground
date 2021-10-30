import React, { useState, useEffect } from 'react';
import { useHistory } from "react-router-dom";
import { Form, Input, Button, Checkbox, PageHeader } from 'antd';
import './Form.css'

const BuildForm = () => {
  const [form] = Form.useForm();
  const history = useHistory();

  const onCheck = async () => {
    try {
      const values = await form.validateFields();
      console.log('Success:', values);
    } catch (errorInfo) {
      console.log('Failed:', errorInfo);
    }
  };
  
  const handleRoute = () =>{ 
    history.push("/");
  }

  return (
    <Form form={form} name="dynamic_rule" className="build-form">
      <PageHeader 
        title="Build Form"
        className="build-form-header"
        onBack={handleRoute}
        // subTitle="TBD"
      />
      <Form.Item
        name="gitRepo"
        label="Git Repository"
        className="build-form-input"
        rules={[
          {
            required: true,
            message: 'This is a required field',
          },
        ]}
      >
        <Input placeholder="Please input the git directory where your project is located" />
      </Form.Item>
      <Form.Item
        name="contextDirectory"
        label="Context Directory"
        className="build-form-input"
        rules={[
          {
            required: false,
          },
        ]}
      >
        <Input placeholder="Please input your context directory" />
      </Form.Item>
      <Form.Item
        name="buildStrategy"
        label="Build Strategy"
        className="build-form-input"
        rules={[
          {
            required: false,
          },
        ]}
      >
        <Input placeholder="Please choose which build strategy you want to use" />
      </Form.Item>
      <Form.Item
        name="dockerfilePath"
        label="Dockerfile Path"
        className="build-form-input"
        rules={[
          {
            required: false,
          },
        ]}
      >
        <Input placeholder="Please input your context directory" />
      </Form.Item>
      <Form.Item
        name="builderImage"
        label="Builder Image"
        className="build-form-input"
        rules={[
          {
            required: false,
          },
        ]}
      >
        <Input placeholder="Please input your builder image" />
      </Form.Item>
      <Form.Item>
        <Button type="primary" className="build-form-submit">
          Build
        </Button>
      </Form.Item>
    </Form>
  );
};

export default BuildForm