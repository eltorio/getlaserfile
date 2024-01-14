# getlaserfile

`getlaserfile` is a Go program designed to be used as a sidecar container for Highcanfly's Gitea Kubernetes Helm chart. It allows public download of specific files from a private Gitea repository.

## License

This program is provided under the GNU Affero General Public License v3.0 (AGPLv3).

## Usage

To use `getlaserfile`, you need to deploy it alongside the Gitea container in your Kubernetes cluster. The program exposes an endpoint that can be used to download specific files from a private Gitea repository.

## Deployment

You can deploy `getlaserfile` using the same Helm chart as your Gitea deployment. Please refer to the Highcanfly's Gitea Kubernetes Helm chart documentation for more details on how to configure and deploy your Helm chart.

## Contributing

Contributions to `getlaserfile` are welcome! Please make sure to read the Contributing Guide before making a pull request.

## Contact

If you have any questions, issues, or feedback, please open an issue in this repository.

## Disclaimer

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
