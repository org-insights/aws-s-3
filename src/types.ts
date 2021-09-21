import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface MyQuery extends DataQuery {
  bucket?: string;
  prefix: string;
}

export const defaultQuery: Partial<MyQuery> = {
  bucket: '',
  prefix: '/',
};

/**
 * These are options configured for each DataSource instance.
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  accessKeyId?: string;
  endpoint?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
  secretAccessKey?: string;
}
