# Web GUI

Start the local scan API:

```sh
kickbucket serve
```

The API listens on `http://127.0.0.1:8080` and exposes:

- `GET /healthz` for readiness checks
- `POST /api/scan` for JSON scans

To run the Next.js planner from this repository, start the website separately:

```sh
cd website
npm run dev
```

Open the local Next.js URL, enter a keyword, select providers, and choose **Run local scan**. The planner previews candidates without the API and displays live results when `kickbucket serve` is running. Use `--wordlist` and `--object-path` with `serve` to change the API configuration.
