// 共通設定
export const config = {
  // API設定
  baseURL: 'http://localhost:8080',
  apiVersion: 'v1',

  // タイムアウト設定
  timeout: '10s',

  // デフォルトヘッダー
  defaultHeaders: {
    'Content-Type': 'application/json',
  },
};

// APIエンドポイントを生成する関数
export function getAPIEndpoint(path) {
  return `${config.baseURL}/api/${config.apiVersion}${path}`;
}

// リクエストパラメータを生成する関数
export function getRequestParams(endpoint, additionalHeaders = {}) {
  return {
    headers: {
      ...config.defaultHeaders,
      ...additionalHeaders,
      'X-Request-ID': `k6-${__VU}-${__ITER}`,
    },
    timeout: config.timeout,
    tags: {
      endpoint: endpoint,
    },
  };
}

// デフォルトのk6オプション
export const defaultOptions = {
  vus: 10,
  duration: '30s',
  thresholds: {
    http_req_duration: ['p(95)<300'], // 95パーセンタイルが300ms未満
    http_req_failed: ['rate<0.01'],   // エラー率1%未満
  },
};