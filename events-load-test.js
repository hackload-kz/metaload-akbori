import http from 'k6/http';
import { sleep, check } from 'k6';

export const options = {

    stages: [
        // Ramp-up to 1000 virtual users over 30 seconds
        { duration: '30s', target: 1000 },
        // Stay at 1000 virtual users for 4 minutes
        { duration: '4m', target: 1000 },
        // Ramp-down to 0 virtual users over 30 seconds
        { duration: '30s', target: 0 },
    ],
};

export default function () {
    const url = `${__ENV.API_URL}/api/events?page=1&pageSize=20`;

    // Capture the response in a variable
    const response = http.get(url);

    // Use a check to formally verify the response
    check(response, {
        'status is 200': (r) => r.status === 200,
    });

    sleep(0.1);
}

// Setup function - runs once before test starts
export function setup() {
    return {
        startTime: Date.now(),
        testVersion: 'v2.1.0',
        environment: 'load_test'
    };
}

import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";
export function handleSummary(data) {
    return {
        'target/summary.html': htmlReport(data),
        'target/summary.json': JSON.stringify(data, null, 2),
    };
}
