import React, { useState, useEffect } from 'react';
import { useHistory } from "react-router-dom";
import { Form, Input, Button, Checkbox, PageHeader } from 'antd';
import { NavBar } from '../NavBar/NavBar.js';
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
    <>
      <NavBar />
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
          <Input placeholder="Link to git repository" />
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
          <Input placeholder="Defaults to root" />
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
          <Input placeholder="Select a Build Strategy" />
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
          <Input placeholder="Defaults to /Dockerfile" />
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
          <Input placeholder="Link to builder image" />
        </Form.Item>
        <Form.Item>
          <Button type="primary" className="build-form-submit">
            Build
          </Button>
        </Form.Item>
      </Form>
    </>
  );
};

export default BuildForm