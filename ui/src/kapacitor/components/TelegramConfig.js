import React, {PropTypes} from 'react';

const TelegramConfig = React.createClass({
  propTypes: {
    config: PropTypes.shape({
      options: PropTypes.shape({
        'chat-id': PropTypes.string.isRequired,
        'disable-notification': PropTypes.bool.isRequired,
        'disable-web-page-preview': PropTypes.bool.isRequired,
        'parse-mode': PropTypes.string.isRequired,
        token: PropTypes.bool.isRequired,
        url: PropTypes.string.isRequired,
      }).isRequired,
    }).isRequired,
    onSave: PropTypes.func.isRequired,
  },

  handleSaveAlert(e) {
    e.preventDefault();

    const properties = {
      chatID: this.chatID.value,
      disableNotification: this.disableNotification.checked,
      disableWebPagePreview: this.disableWebPagePreview.checked,
      parseMode: this.parseMode.checked,
      token: this.token.value,
      url: this.url.value,
    };

    this.props.onSave(properties);
  },

  render() {
    const {options} = this.props.config;
    const {url, token} = options;
    const chatID = options['chat-id'];
    const disableNotification = options['chat-id'];
    const disableWebPagePreview = options['disable-web-page-preview'];
    const parseMode = options['parse-mode'];

    return (
      <div>
        <h4 className="text-center">Telegram Alert</h4>
        <br/>
        <p>You can have alerts sent to Telegram by entering info below.</p>
        <form onSubmit={this.handleSaveAlert}>
          <div className="form-group col-xs-12">
            <label htmlFor="url">Telegram URL</label>
            <input className="form-control" id="url" type="text" ref={(r) => this.url = r} defaultValue={url || ''}></input>
          </div>

          <div className="form-group col-xs-12">
            <label htmlFor="token">Token</label>
            <input className="form-control" id="token" type="text" ref={(r) => this.token = r} defaultValue={token || ''}></input>
            <label className="form-helper">Note: a value of <code>true</code> indicates the Telegram token has been set</label>
          </div>

          <div className="form-group col-xs-12">
            <label htmlFor="chat-id">Chat ID</label>
            <input className="form-control" id="chat-id" type="text" ref={(r) => this.chatID = r} defaultValue={chatID || ''}></input>
          </div>

          <div className="form-group col-xs-12">
            <div className="form-control-static">
              <input id="enableParseMode" type="checkbox" defaultChecked={parseMode} ref={(r) => this.parseMode = r} />
              <label htmlFor="enableParseMode">Enable Parse Mode</label>
            </div>
          </div>

          <div className="form-group col-xs-12">
            <div className="form-control-static">
              <input id="disableWebPagePreview" type="checkbox" defaultChecked={disableWebPagePreview} ref={(r) => this.disableWebPagePreview = r} />
              <label htmlFor="disableWebPagePreview">Disable Web Page Preview</label>
            </div>
          </div>

          <div className="form-group col-xs-12">
            <div className="form-control-static">
              <input id="disableNotification" type="checkbox" defaultChecked={disableNotification} ref={(r) => this.disableNotification = r} />
              <label htmlFor="disableNotification">Disable Notification</label>
            </div>
          </div>

          <div className="form-group form-group-submit col-xs-12 col-sm-6 col-sm-offset-3">
            <button className="btn btn-block btn-primary" type="submit">Save</button>
          </div>
        </form>
      </div>
    );
  },
});

export default TelegramConfig;
