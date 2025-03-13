import { check } from "k6";
import http from "k6/http";

const BASE_URL = __ENV.BASE_URL || "http://localhost:8000";

export const options = {
    stages: [
        { duration: "2m", target: 20 },
        { duration: "2m", target: 20 },
        { duration: "2m", target: 50 },
        { duration: "5m", target: 50 },
        { duration: "2m", target: 120 },
        { duration: "5m", target: 120 },
        { duration: "2m", target: 200 },
        { duration: "5m", target: 200 },
        { duration: "5m", target: 0 },
    ],
    summaryTrendStats: ["min", "med", "avg", "count", "max", "p(95)", "p(99)"],
    thresholds: {
        http_req_duration: ["p(95)<500", "p(99)<700"],
    },
};

export default function () {
    const res = http.get(`${BASE_URL}/`);

    check(res, {
        "status is 200": (r) => r.status === 200,
    });
}
