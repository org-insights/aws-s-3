import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { MyDataSourceOptions, MySecureJsonData } from './types';

const { SecretFormField, FormField } = LegacyForms;

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions> {}

interface State {}

export class ConfigEditor extends PureComponent<Props, State> {
  onAccessKeyIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      accessKeyId: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  onEndpointChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      endpoint: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  // Secure field (only sent to the backend)
  onSecretAccessKeyChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonData: {
        secretAccessKey: event.target.value,
      },
    });
  };

  onResetSecretAccessKey = () => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        secretAccessKey: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        secretAccessKey: '',
      },
    });
  };

  render() {
    const { options } = this.props;
    const { jsonData, secureJsonFields } = options;
    const secureJsonData = (options.secureJsonData || {}) as MySecureJsonData;

    return (
      <div className="gf-form-group">
        <div className="gf-form">
          <FormField
            label="Access Key ID"
            labelWidth={8}
            inputWidth={20}
            onChange={this.onAccessKeyIdChange}
            value={jsonData.accessKeyId || ''}
            placeholder="Access Key ID"
          />
        </div>

        <div className="gf-form-inline">
          <div className="gf-form">
            <SecretFormField
              isConfigured={(secureJsonFields && secureJsonFields.apiKey) as boolean}
              value={secureJsonData.secretAccessKey || ''}
              label="Secret Access Key"
              placeholder="Enter Secret Access Key"
              labelWidth={8}
              inputWidth={20}
              onReset={this.onResetSecretAccessKey}
              onChange={this.onSecretAccessKeyChange}
            />
          </div>
        </div>
        <div className="gf-form">
          <FormField
            label="Endpoint"
            labelWidth={8}
            inputWidth={20}
            onChange={this.onEndpointChange}
            value={jsonData.endpoint || ''}
            placeholder="URL endpoint to make API calls to"
          />
        </div>
      </div>
    );
  }
}
