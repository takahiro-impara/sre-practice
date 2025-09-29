import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';
import { config, getAPIEndpoint, getRequestParams, defaultOptions } from '../config/common.js';

const listUsersRate = new Rate('list_users_rate');
const listUsersTrend = new Trend('list_users_trend');

export const options = defaultOptions;

export function setup() {
  const healthCheck = http.get(`${config.baseURL}/healthz`);
  check(healthCheck, { 'Service is healthy': (r) => r.status === 200 });

  console.log('Service is healthy');

  return { startTime: Date.now() };
}

// メインテスト関数
export default function (data) {
  // ページネーションパラメータをランダムに変更
  const limit = Math.floor(Math.random() * 50) + 10; // 10-59
  const offset = Math.floor(Math.random() * 100); // 0-99

  const params = getRequestParams('list-users');

  // リクエスト実行とレスポンス時間測定
  const startTime = new Date();
  const response = http.get(`${getAPIEndpoint('/users')}?limit=${limit}&offset=${offset}`, params);
  const duration = new Date() - startTime;

  listUsersRate.add(1);
  listUsersTrend.add(duration);

  // レスポンス検証
  const checkResult = check(response, {
    'Status is 200': (r) => r.status === 200,
    'Response has users array': (r) => {
      try {
        const body = JSON.parse(r.body);
        return Array.isArray(body.users);
      } catch (e) {
        return false;
      }
    },
    'Response time < 300ms': (r) => r.timings.duration < 300
  });

  // エラーがあればログ出力
  if (!checkResult) {
    console.error('Test failed:', response.status, response.body);
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
    test_type: 'user_list',
    duration: duration,
    timestamp: new Date().toISOString(),
  };

  http.post('http://metrics-collector:8080/results', JSON.stringify(summary));
}
