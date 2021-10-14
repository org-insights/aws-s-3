import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface MyQuery extends DataQuery {
  bucket?: string;
  prefix: string;
  metric: number;
}

export const defaultQuery: Partial<MyQuery> = {
  bucket: '',
  prefix: '/',
  metric: 0,
};

/**
 * These are options configured for each DataSource instance.
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  authenticationProvider: number;
  accessKeyId?: string;
  endpoint?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
  secretAccessKey?: string;
}
