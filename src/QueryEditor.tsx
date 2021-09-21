import { defaults } from 'lodash';

import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, MyDataSourceOptions, MyQuery } from './types';

const { FormField } = LegacyForms;

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  onBucketChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, bucket: event.target.value });
  };

  onPrefixChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, prefix: event.target.value });
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { bucket, prefix } = query;

    return (
      <div className="gf-form">
        <FormField width={4} value={bucket} onChange={this.onBucketChange} label="Bucket" tooltip="Bucket name" />
        <FormField
          width={10}
          value={prefix || ''}
          onChange={this.onPrefixChange}
          label="Prefix"
          tooltip="Prefix path in bucket"
        />
      </div>
    );
  }
}
