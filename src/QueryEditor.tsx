import { defaults } from 'lodash';

import React, { ChangeEvent, PureComponent } from 'react';
import { InlineField, Input, LegacyForms, Select } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, MyDataSourceOptions, MyQuery } from './types';

const { FormField } = LegacyForms;

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

const metricOptions = [
  { label: 'Size', value: 0, description: 'Keys size in bytes' },
  { label: 'Number of keys', value: 1, description: 'Number of keys' },
];

export class QueryEditor extends PureComponent<Props> {
  onBucketChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, bucket: event.target.value });
  };

  onPrefixChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, prefix: event.target.value });
  };

  onMetricChange = (event: SelectableValue<number>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, metric: event.value || 0 });
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { bucket, prefix, metric } = query;

    return (
      <div className="gf-form">
        <FormField width={4} value={bucket} onChange={this.onBucketChange} label="Bucket" tooltip="Bucket name" />
        <InlineField label="Prefix" tooltip="Prefix path in bucket" grow>
          <Input placeholder="Inline input" css={undefined} value={prefix || ''} onChange={this.onPrefixChange} />
        </InlineField>
        <InlineField label="Metric" labelWidth={10}>
          <Select options={metricOptions} width={20} value={metric} onChange={this.onMetricChange} />
        </InlineField>
      </div>
    );
  }
}
