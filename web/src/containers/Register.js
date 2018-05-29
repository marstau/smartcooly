import { ResetError } from '../actions';
import { UserRegister } from '../actions/user';
import React from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';
import { Button, Form, Input, Icon, Tooltip, notification } from 'antd';

class Register extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      windowHeight: window.innerHeight || 720,
      messageErrorKey: '',
    };

    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    const { dispatch } = this.props;
    const { messageErrorKey } = this.state;
    const { user } = nextProps;

    if (!messageErrorKey && user.message) {
      this.setState({
        messageErrorKey: 'userRegisterError',
      });
      notification['error']({
        key: 'userRegisterError',
        message: 'Error',
        description: String(user.message),
        onClose: () => {
          if (this.state.messageErrorKey) {
            this.setState({ messageErrorKey: '' });
          }
          dispatch(ResetError());
        },
      });
    }

    if (!user.loading && user.token) {
      if (this.state.messageErrorKey) {
        this.setState({ messageErrorKey: '' });
        notification.close(this.state.messageErrorKey);
      }
      browserHistory.push('/');
    }
  }

  componentWillMount() {
    this.order = 'id';
  }

  componentWillUnmount() {
    const { dispatch } = this.props;

    dispatch(ResetError());
  }

  Login(e) {
    browserHistory.push('/login');
  }

  handleSubmit(e) {
    const { form, dispatch } = this.props;

    if (e) {
      e.preventDefault();
    }

    form.validateFields((errors, values) => {
      if (errors) {
        return;
      }

      const req = {
        id: 0,
        username: values.username,
        level: values.level,
        email: values.email,
      };

      dispatch(UserRegister(values.cluster, req, values.password));
    });
  }

  render() {
    const { windowHeight } = this.state;
    const { getFieldDecorator } = this.props.form;
    const formItemLayout = {
      wrapperCol: { offset: 9, span: 6 },
    };
    const cluster = localStorage.getItem('cluster') || document.URL.slice(0, -9);

    return (
      <div style={{ paddingTop: windowHeight > 600 ? (windowHeight - 500) / 2 : windowHeight > 400 ? (windowHeight - 350) / 2 : 25 }}>
        <h1 style={{
          margin: 24,
          fontSize: '30px',
          textAlign: 'center',
        }}>Smartcooly</h1>
        <Form horizontal onSubmit={this.handleSubmit} onClick={this.handleClick}>
          <Form.Item
            {...formItemLayout}
          >
            <Tooltip placement="right" title="Cluster URL">
              {getFieldDecorator('cluster', {
                rules: [{ type: 'url', required: true }],
                initialValue: cluster,
              })(
                <Input addonBefore={<Icon type="link" />} placeholder="http://127.0.0.1:9876" />
              )}
            </Tooltip>
          </Form.Item>
          <Form.Item
            {...formItemLayout}
          >
            <Tooltip placement="right" title="Username">
              {getFieldDecorator('username', {
                rules: [{ required: true }],
              })(
                <Input addonBefore={<Icon type="user" />} placeholder="username" />
              )}
            </Tooltip>
          </Form.Item>
          <Form.Item
            {...formItemLayout}
          >
            <Tooltip placement="right" title="Password">
              {getFieldDecorator('password', {
                rules: [{ required: true }],
              })(
                <Input addonBefore={<Icon type="lock" />} type="password" placeholder="password" />
              )}
            </Tooltip>
          </Form.Item>
          <Form.Item
            {...formItemLayout}
          >
            <Tooltip placement="right" title="Email">
              {getFieldDecorator('email', {
                rules: [{ required: true }],
              })(
                <Input addonBefore={<Icon type="mail" />} placeholder="email" />
              )}
            </Tooltip>
          </Form.Item>
          <Form.Item wrapperCol={{ span: 6, offset: 9 }} style={{ marginTop: 24 }}>
            <Button type="primary" htmlType="submit" className="login-form-button">Register</Button>
          </Form.Item>
          <a href="#" onClick={this.Login}> <h5 style= {{ margin: 24, textAlign: 'center', }}>Login</h5></a>
        </Form>
      </div>
    );
  }
}

const mapStateToProps = (state) => ({
  user: state.user,
});

export default Form.create()(connect(mapStateToProps)(Register));
