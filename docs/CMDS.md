# Commands

## Global Commands

You can use all commands from everywhere in your project. 
The tool automatically detect the context.

  - **[project](#project)** : Manage projects
  - **[component](#component)** : Manage Kubernetes components
  - **[cluster](#cluster)** : Manage Local Kubernetes clusters
  - **[check](#check)**:Check the required tools
  - **[completion](../README.md#shell-completion)** : Output shell completion code for the specified shell
  - **help** : Help about any command
  - **version** : Print the version information

<a id="global-flag"></a>

## Global Flags

  **-c, --config** *string* : Project configuration file

  **-h, --help** : help display
  
  **-v, --verbose** : Enable verbose logging

<a id="project"></a>

## Project Command

All commands to manage project.

  - **[create](#project-create)** : Create a project configuration
  - **[edit](#project-edit)** : Edit a project configuration
  - **[install](#project-install)** : Install a project

<a id="project-create"></a>

### Project create

Create standard `gikops.yaml` file in your folder and install all
specified configurations.

```shell
gikopsctl project create

## Alias usage
gop create
```

Follow instructions:

```shell
What is the name for the project?
# default current directory name.

Do you want to create only standard kind cluster ?[y/N]
# If you select Yes, only a local Kind cluster should be created.

# If you select No, you should specify all clusters you need.
Select cluster type:
- kind
- basic

What is the name for the cluster?
# Cluster name. It correspond to kubernetes registred context in your
# Kubectl.

# If you select Kind cluster you must select the container provider
# Becareful, select only provider you have installed on your host.
Select container runtime provider:
- docker
- podman
- nerdctl

Do you want to add another cluster? [y/N]
# Loop to add other clusters


Do you want to add components? [y/N]
# Only if you want preset internal core components.
# It's usefull to start a boilerplate with basic components.

Select [type] components to install:
> [x] component 1
> [x] component 2
> [x] component 3 
# Select each core component you want to add to your project

Do you want to add components folders? [y/N]
# You can create other components namespaces as you whish
```

Result:

```shell
YourFolder
├── clusters
│   ├── ...
│   └── local
│       ├── kind.yaml
│       └── overrides
├── core
│   ├── cert-manager
│   │   └── ...
│   ├── mkcert
│   │   └── ...
│   ├── prometheus
│   │   └── ...
│   └── traefik
│       └── ...
├── gikops.yaml
└── project
```

<a id="project-install"></a>

### Project install

If you have only the `gikops.yaml` file, this command generate all other folders and files, folling the configuration specified inside.

Only use it if you have **only** the `gikops.yaml` file.

```shell
gikopsctl project install

## Alias usage
gop install
```

<a id="project-edit"></a>

### Project edit

Edit project configurations through command line.

- **[name](#project-edit-name)** : Set the name of the project
- **[add](#project-edit-add)** : Add parameters to the project
- **[delete](#project-edit-delete)**: : Delete parameters from the project

<a id="project-edit-name"></a>

#### Project edit name

Set the name of the project

```shell
gikopsctl project edit name <new_name>

## Alias usage
gop edit name <new_name>
```

<a id="project-edit-add"></a>

#### Project edit add

Adding parameters to project configuration.

- **[cluster](project-edit-add-cluster)** : Add a cluster to the project
- **[component](project-edit-add-component)** : Add components namespace to the project

<a id="project-edit-add-cluster"></a>

##### Project edit add cluster

Add new cluster to your project.
Same command result as [Cluster add](#cluster-add).

```shell
gikopsctl project edit add cluster

## Alias usage
gop edit add cluster
```

<a id="project-edit-add-component"></a>

##### Project edit add component

Add components namespace to the project.
It also create a new sub folder into your project root.

```shell
gikopsctl project edit add component <componens_namespace_name>

## Alias usage
gop edit add component <componens_namespace_name>
```

<a id="project-edit-delete"></a>

#### Project edit delete

<div class="alert">
  <strong>Warning!</strong> Cluster removal don't work from this commande
  dont work yet. Please use Clusters commands.
</div>

cluster     Remove cluster from the project

<a id="component"></a>

## Component Command

- **[create](#component-create)** : Create a component
- **[check](#component-check)** : Check a component configuration
- **[init](#component-init)** : Initialize a component
- **[apply](#component-apply)** : Apply a component to the cluster
- **[delete](#component-delete)** : Delete a component

<a id="component-flags"></a>

### Component flags

**-a, --all** : Select all components of the targeted namespace

**--folder** *string* Namespace targeted. If your are into a namespace folder
this value is auto-completed.

[Global flags](#global-flags)

<a id="component-create"></a>

### Component create

Create new component in your targeted components namespace.

```shell
gikopsctl component create [--folder name]

## Alias usage
goco create [--folder name]
```

Follow instructions:

```shell
What is the name for the component?
# Set tne name of the new component

Specific namespace of the component
# You can override kubernetes Namespace

Do you want to use helm? [y/N]
# Use this if your component com from an Helm repository

Helm repository name?
# Set Helm repository name

Helm repository URL?
# Set Helm repository URL

Helm chart name?
# Set Helm chart name

Specifique version?
# Set chart specific version

Different CRDs chart? [y/N]
# Use this option if CRDs are defined in an other Chart

Do you want to use kustomize? [y/N]
# Use this option if you want to download static files to implement
# inside the kustomize file.
```

<a id="component-check"></a>

### Component check

Check if component configuration is correct.

```shell
gikopsctl component check [--folder name]

## Alias usage
goco check [--folder name]
```

<a id="component-init"></a>

### Component init

Generate necessary files from config definitions.
If Helm was set into configuration file, init make:
- Update helm repo (or add)
- Get default values (if not already exist)
- Generate Helm template result (temporary)
- Cut into multiple files
- Create Kustomize yaml file (if not exist)

```shell
gikopsctl component init [component-name] [flags]

## Alias usage
goco init [component-name] [flags]
```

If you are into the component folder `component-name` is not
necessary.

#### Flags

**-k, --keep-tmp** : Keep the temporary Helm chart files.

**-o, --only** : Only initialize the component in the current folder, no dependecies.

[Component flags](#component-flags)

<a id="component-apply"></a>

### Component apply

Apply CRDs in the first step.
Generate computed yaml file with all manifests corresponding to 
selected cluster using Kustomize.
Apply computed file into selected cluster.

```shell
gikopsctl component apply [component-name] [flags]

## Alias usage
goco apply [component-name] [flags]
```

If you are into the component folder `component-name` is not
necessary.

#### Flags

**-e, --env** *string* : Cluster to target (default "local").

**-m, --mode** *string* : Mode to apply (all, crds, manifests) (default "all").

**-o, --only** : Only the component in the current folder without dependencies.

**-b, --only-build** : Only build step to generate computed.yaml file/

[Component flags](#component-flags)

<a id="component-delete"></a>

### Component delete

Delete from selected cluster all ressources corresponding to the 
component.

```shell
gikopsctl component delete [component-name] [flags]

## Alias usage
goco delete [component-name] [flags]
```

If you are into the component folder `component-name` is not
necessary.

#### Flags

**-e, --env** *string* : Cluster to target (default "local").

**-m, --mode** *string* : Mode to apply (all, crds, manifests) (default "all").

**--force** : Force delete component

[Component flags](#component-flags)

<a id="cluster"></a>

## Cluster Command

  - **[add](#cluster-add)** : Add a cluster to the project
  - **[delete](#cluster-delete)** : Delete a cluster from the project
  - **[install](#cluster-install)** : Install the Kubernetes environment and components
  - **[uninstall](#cluster-uninstall)** : Uninstall the Kubernetes environment


<a id="cluster-add"></a>

## Cluster add

Add a new cluster into the project.
It generate a new cluster folder into all components.

```shell
gikopsctl cluster add [flags]

## Alias usage
gocu add [flags]
```

follow instructions.

<a id="cluster-delete"></a>

## Cluster delete

Delete selected cluster and all component configs.

```shell
gikopsctl cluster delete [cluster-name] [flags]

## Alias usage
gocu delete [cluster-name] [flags]
```

<a id="cluster-install"></a>

## Cluster install

Install and start selected cluster.

```shell
gikopsctl cluster install [cluster-name] [flags]

## Alias usage
gocu install [cluster-name] [flags]
```

<a id="cluster-uninstall"></a>

## Cluster uninstall

Stop and remove selected cluster.

```shell
gikopsctl cluster uninstall [cluster-name] [flags]

## Alias usage
gocu uninstall [cluster-name] [flags]
```

<a id="check"></a>

## Check Command

TODO
