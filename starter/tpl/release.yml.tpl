# .github/workflows/release.yml
name: goreleaser

on:
  push:
    # run only against tags
    tags:
      - "*"

permissions:
  contents: write
  # packages: write
  # issues: write
  # id-token: write

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: stable
      # More assembly might be required: Docker logins, GPG, etc.
      # It all depends on your needs.
      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v6
        with:
          # either 'goreleaser' (default) or 'goreleaser-pro'
          distribution: goreleaser-pro
          # 'latest', 'nightly', or a semver
          version: "~> v2"
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          # Your GoReleaser Pro key, if you are using the 'goreleaser-pro' distribution
          GORELEASER_KEY: ${{ secrets.GORELEASER_KEY }}
          # Uncomment for announcements
          # BLUESKY_APP_PASSWORD: ${{ secrets.BLUESKY_APP_PASSWORD }}