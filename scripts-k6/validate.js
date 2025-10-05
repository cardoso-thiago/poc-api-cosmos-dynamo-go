import http from 'k6/http';
import { SharedArray } from 'k6/data';

export const options = {
  scenarios: {
    ramping_load: {
      executor: 'ramping-arrival-rate',
      startRate: 1,
      timeUnit: '1s',
      preAllocatedVUs: 10,
      maxVUs: 100,
      stages: [
        { target: 100, duration: '10s' },
        { target: 100, duration: '1m' },
      ],
    },
  },
};

const uuids = new SharedArray('uuids', function() {
    return open('../docker/uuids.txt').split('\n').filter(u => u.trim().length > 0);
});

export default function () {
    const uuid = uuids[Math.floor(Math.random() * uuids.length)];
    const res = http.get(`http://localhost:8888/items/${uuid}`);
}