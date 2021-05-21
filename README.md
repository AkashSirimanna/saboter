# saboter

Saboter is a simple, yet practical application of Chaos Engineering in Kubernetes.
It sabotages prod... On purpose

## Running Saboter

Use the Dockerfile

OR

```
go build -o saboter
./saboter --kubepath <path> --exclude <path> --rate=<num> --interval=<num>
```

# Args

- kubepath: Path to your Kubernetes config, defaults to ~/.kube/config
- exclude: Path to a file containing days (YYYY-MM-DD) to exclude running saboter on
- rate: How many pods to kill every interval
- interval: Interval at which to kill pods

Example

```
./saboter --kubepath ~/.kube/config --exclude ~/code/days --rate=2 --interval=1
```

Example of exclude file

```
2021-05-10
2021-05-11
```

## Contributing
Any and all PRs are welcome.
