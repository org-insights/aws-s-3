import React, { ChangeEvent, PureComponent } from 'react';
import { InlineField, LegacyForms, Select } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps, SelectableValue } from '@grafana/data';
import { MyDataSourceOptions, MySecureJsonData } from './types';

const { SecretFormField, FormField } = LegacyForms;

const authOptions = [
  { label: 'AWS SDK Default', value: 0, description: 'Authentication with Assume Role' },
  { label: 'Acess & Secret Keys', value: 1, description: 'Authentication with Access Key ID and Secret Access Key' },
];

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions> {}

interface State {}

export class ConfigEditor extends PureComponent<Props, State> {
  onAuthenticationProviderChange = (event: SelectableValue<number>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      authenticationProvider: event.value || 0,
    };
    console.log(event.value);
    onOptionsChange({ ...options, jsonData });
  };

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
          <InlineField label="Authentication Provider" labelWidth={20}>
            <Select
              options={authOptions}
              width={40}
              value={jsonData.authenticationProvider}
              allowCustomValue
              onChange={this.onAuthenticationProviderChange}
            />
          </InlineField>
        </div>

        {jsonData.authenticationProvider === 1
          ? [
              <div className="gf-form">
                <FormField
                  label="Access Key ID"
                  labelWidth={10}
                  inputWidth={20}
                  onChange={this.onAccessKeyIdChange}
                  value={jsonData.accessKeyId || ''}
                  placeholder="Access Key ID"
                />
              </div>,

              <div className="gf-form-inline">
                <div className="gf-form">
                  <SecretFormField
                    isConfigured={(secureJsonFields && secureJsonFields.apiKey) as boolean}
                    value={secureJsonData.secretAccessKey || ''}
                    label="Secret Access Key"
                    placeholder="Enter Secret Access Key"
                    labelWidth={10}
                    inputWidth={20}
                    onReset={this.onResetSecretAccessKey}
                    onChange={this.onSecretAccessKeyChange}
                  />
                </div>
              </div>,
            ]
          : null}

        <div className="gf-form">
          <FormField
            label="Endpoint"
            labelWidth={10}
            inputWidth={20}
            onChange={this.onEndpointChange}
            value={jsonData.endpoint || ''}
            placeholder="Optionally, specify a custom endpoint for S3"
          />
        </div>
      </div>
    );
  }
}
