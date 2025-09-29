import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';
import { config, getAPIEndpoint, getRequestParams, defaultOptions } from '../config/common.js';

const createUserRate = new Rate('create_user_rate');
const createUserTrend = new Trend('create_user_trend');

export const options = defaultOptions;

export function setup() {
  const healthCheck = http.get(`${config.baseURL}/healthz`);
  check(healthCheck, { 'Service is healthy': (r) => r.status === 200 });

  console.log('Service is healthy');

  return { startTime: Date.now() };
}

// メインテスト関数
export default function () {
  const email = `test-${Math.floor(Math.random() * 900000) + 100000}@example.com`;
  const name = `Test User ${Math.floor(Math.random() * 900000) + 100000}`;
  const password = 'Password123';
  const payload = {
    email: email,
    name: name,
    password: password,
  };

  const params = getRequestParams('create-user');

  // リクエスト実行とレスポンス時間測定
  const startTime = new Date();
  const response = http.post(getAPIEndpoint('/users'), JSON.stringify(payload), params);
  const duration = new Date() - startTime;

  createUserRate.add(1);
  createUserTrend.add(duration);

  // レスポンス検証
  const checkResult = check(response, {
    'Status is 201': (r) => r.status === 201,
    // 'Response has user ID': (r) => r.body && r.body.id,
    'Response time < 300ms': (r) => r.timings.duration < 300
  });

  // エラーがあればログ出力
  if (!checkResult) {
    console.error('Test failed:', response);
  }

  // 思考時間のシミュレーション（1-3秒のランダム待機）
  sleep(Math.random() * 2 + 1);
}

// ティアダウン関数（テスト後に1回実行）
export function teardown(data) {
  const duration = (Date.now() - data.startTime) / 1000;
  console.log(`Test completed in ${duration} seconds`);

  // 結果のサマリー送信（外部システムへ）
  const summary = {
    test_type: 'user_creation',
    duration: duration,
    timestamp: new Date().toISOString(),
  };

  http.post('http://metrics-collector:8080/results', JSON.stringify(summary));
}
