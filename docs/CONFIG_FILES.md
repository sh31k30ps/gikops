# Configurations files

## Config files

- __[📦 gikops.yaml](#gikops-yaml)__ : Global project configuration file.
- __[🧩 gikcpnt.yaml](#gikcpnt-yaml)__ : Componant configuration file.

<a id="gikops-yaml"></a>

## gikops.yaml

This file must be in root project folder. It contain information about
componants namespaces and clusters.

```yaml
kind: Project
apiVersion: gikopsctl.config.k8s.io/v1alpha1
metadata:
  name: string # Name of the project
components: [] # Components namespaces List
clusters: []   # Clusters definitions List
```

### Project

**components** (*ArrayOf[[ProjectComponent](#ProjectComponent)]*) :<br/>
Components is a list of folders containing component configurations.

**clusters** (*ArrayOf[[ProjectCluster](#ProjectCluster)]*) :<br/>
Clusters is a list of clusters and there configurations.

<a id="ProjectComponent"></a>

### Project Component

```yaml
  - name: core
    require:
      - ingress/traefik
      - monitoring/prometheus
      - security/cert-manager
      - security/mkcert
  - name: project
```

**name** *(string)* mandatory : <br/>
Name of folder containing components configurations. The folder must exist 
into your project.

**require** *(ArrayOf[string])* : <br/>
List of internal preset components that are required by the components namespace.

Existing presets components:

| Context    | Name of componnent | Full path             |
|------------|--------------------|-----------------------|
| Ingress    | Traefik            | ingress/traefik       |
| Monitoring | Prometheus         | monitoring/prometheus |
| Security   | Cert Manager       | security/cert-manager |
| Security   | MKCert             | security/mkcert       | 

<a id="ProjectCluster"></a>

### Project Cluster 

```yaml
name: local
kindConfig: {} # Only if is an cluster using local kind cluster
```

**name** *(string)* mandatory : <br/>
Name of the cluster

**kindConfig** *([KindCluster](#KindCluster))* : <br/>
Only use this configuration if you want a local Kind Cluster for
developpment.

<a id="KindCluster"></a>

### Kind Cluster configuration

```yaml
configFile: string
clusterName: string
overridesFolder: [] # []string
provider: string
```

**clusterName** *(string)* (default: empty) : <br/>
Intenal name of the local Kind cluster. Only use it if you want
a different name of Cluster name. It could be useful if you have multiple 
local Kind Cluster.

**configFile** *(string)* (default: "kind.yaml") : <br/>
Name of Kind configuration file. 

**overridesFolder** *(ArrayOf[string])* (default: ["overrides"]) : <br/>
List of folders containing specific overrides configurations for
Kind cluster.

**provider** *(string)* (default: "docker") : <br/>
Provider for the local Kind cluster.

List of available providers: 
| Provider |
|----------|
|  docker  |
|  podman  |
|  nerdctl |


<a id="gikcpnt-yaml"></a>

## gikcpnt.yaml

This file must be present into each component folders.
It contains all configuration for the component.

```yaml
kind: Component
apiVersion: gikopsctl.config.k8s.io/v1alpha1
metadata:
  name: string  # Name of the component
helm: {}        # Helm config only if is based on. 
files: {}
exec: {}
dependsOn: []
clusters: []
```

### Component

**helm** *([HelmConfig](#helm-config))* : <br/> 
Helm chart configuration for this component.
If unset, the component will not use Helm for deployment.

**kustomize** *([KustomizeConfig](#kustomize-config))* : <br/>
Kustomize configuration for this component.
If unset, the component will not use Kustomize for deployment.

**files** *([ComponentFiles](#component-files))* : <br/>
File definitions for this component.
If unset, the component will not use direct file management

**exec** *([ComponentExec](#component-exec))* : <br/>
Exec contains commands to execute before or after component deployment.
If unset, no commands will be executed.

**dependsOn** *(ArrayOf[string])* : <br/>
List of components that must be deployed before this component

**clusters** *(ArrayOf[string])* : <br/>
List of clusters where the component will be deployed.
by default is all clusters. Some structural componants could be 
deployed only in production.

<a id="helm-config"></a>

### Helm Config

```yaml
repo: string
repo-url: string
version: string
chart: string
crds-chart: string
crds-version: string
before: {}
after: {}
```

**repo** *(string)* mandatory :<br/>
Name of the Helm repository.

**repo-url** *(string)* mandatory :<br/>
URL of the Helm repository.

**version** *(string)* :<br/>
Version of the Helm chart to use. By default Helm use `latest`.

**chart** *(string)* mandatory :<br/>
Chart name used by Helm.

**crds-chart** *(string)* :<br/>
Name of a separate chart containing CRDs.
If unset, CRDs are assumed to be in the main chart.

**crds-version** *(string)* :<br/>
Version of the CRD chart.
If unset and CRDChart is set, Version will be used.

**before** *([HelmBeforeInitConfig](#helm-before))* :<br/>
Contains actions to perform before chart installation

**after** *([HelmAfterInitConfig](#helm-after))* :<br/>
Contains actions to perform after chart installation

<a id="helm-before"></a>

### Helm Before Init Config

```yaml
uploads: []
```

**uploads** *(ArrayOf[[Upload](#exec-upload)])* :<br/>
Specifies files to be uploaded before installation

<a id="helm-after"></a>

### Helm After Init Config

```yaml
uploads: []
resolves: []
renames: []
concats: []
```

**uploads** *(ArrayOf[[Upload](#exec-upload)])* :<br/>
Specifies files to be uploaded before installation

**resolves** *(ArrayOf[string])* :<br/>
Targeted files if they contains only an URL

**renames** *(ArrayOf[[Rename](#exec-rename)])* :<br/>
Specifies files to be renamed after installation

**concats** *(ArrayOf[[Concat](#exec-concat)])* :<br/>
Specifies files to be concatenated after installation


<a id="exec-upload"></a>

### Upload

```yaml
name: string
url: string
```

**name** *(string)* mandatory :<br/>
Target name for the uploaded file

**url** *(string)* mandatory :<br/>
Source URL to download from

<a id="exec-rename"></a>

### Rename

```yaml
original: string
renamed: string
```

**original** *(string)* mandatory :<br/>
Original file name.

**renamed** *(string)* mandatory :<br/>
New file name

<a id="exec-concat"></a>

### Concat

```yaml
folder: string
includes: []
output: string
```

**folder** *(string)* mandatory :<br/>
Directory containing source files.

**includes** *(ArrayOf[string])* :<br/>
Pattern matching files to include.

**output** *(string)* mandatory :<br/>
Target concatenated file.

<a id="kustomize-config"></a>

### Kustomize Config

```yaml
urls: []
```

**urls** *(ArrayOf[string])* :<br/>
List of URLs to download into the component as base of it.

<a id="component-files"></a>

### Component Files

```yaml
crds: string
skipCRDs: boolan
keep: []
```
**crds** *(string)* :<br/>
Path to custom resource definitions. It could be a file or a folder.
By default,it's target `crds.yaml`.

**skipCRDs** *(boolean)* :<br/>
When the project don't have CRDs, you must set this value to `true`.
usually it is used for your own components.

**keep** *(ArrayOf[string])* :<br/>
List of files to preserve during init operations.
By default only `kustomization.yaml` is preserved.
You can add other files as you wish.

<a id="component-exec"></a>

### Component Exec

```yaml
before: []
after: []
```

**before** *(ArrayOf[string])* :<br/>
List of shell commands to execute locally before deployment

**after** *(ArrayOf[string])* :<br/>
List of shell commands to execute locally after deployment success