# Build Spec File Example using Java and Maven

This is an example Hello World project using JDK 15 and Maven. With the [OCI DevOps service](https://www.oracle.com/devops/devops-service/) 
and this project, you'll be able to build this application as a JAR file, create a Docker image from it and store it in
[OCI Container Registry.](https://docs.oracle.com/en-us/iaas/Content/Registry/Concepts/registryoverview.htm)

In this example, you'll install JDK 15, build a JAR file of this Java Hello World Spring Boot Web app, create a Docker 
image from it and store it in the OCI Container Registry, all using the OCI DevOps service!

Let's go!

## Building the application locally

### Download the repo
The first step to get started is to download the repository to your local workspace

```shell
git clone git@github.com:dlcbld/BuildSpecSampleFiles.git
cd java-maven
```

### Install and run the application

1. Install JDK 15 on your system: https://download.java.net/openjdk/jdk15/ri/openjdk-15+36_linux-x64_bin.tar.gz
2. Compile the app:

   ```mvn clean install```
3. To verify, make sure the JAR file is created in your target folder.

### Build a container image for the app
You can locally build a container image using docker (or your favorite container image builder), to verify that you can run the app within a container.

```
docker build -t hello-world:1.0 .
```

Verify that your image was built, with `docker images`

## Building the application in OCI DevOps
Now that you've seen you can locally build this app, let's try this out through OCI DevOps Build service.

### Create External Connection to your Git repository

1. Create a [DevOps Project](https://docs.oracle.com/en-us/iaas/Content/devops/using/devops_projects.htm) or use and an existing project.
2. In your DevOps project, create an External Connection to your GitHub repository which holds your application.
   - Create a Personal Access Token (PAT): https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token
   - In the OCI Console, Go to Identity & Security -> Vault and create a [Vault]( https://docs.oracle.com/en-us/iaas/Content/KeyManagement/Concepts/keyoverview.htm) in compartment of your own choice.
   - Create a Master Key that will be used to encrypt the PATs.
   - Select Secrets from under Resources and create a secret using PAT obtained from GitHub account.
   - Make a note of the OCID of the secret.
   - Now, go to the desired project and select External Connection from the resources.
   - Select type as GitHub and provide OCID of the secret under Personal Access Token.
   - Finally, allow BuildPipeline (dynamic group with DevOps Resources) to use a PAT secret by writing a policy in the root compartment as: ``` Allow dynamic-group dg-with-devops-resources to manage secret-family in tenancy```

### Setup your Build Pipeline
Create a new Build Pipeline to build, test and deliver artifacts from your GitHub Repository.

### Managed Build stage
In your Build Pipeline, first add a Managed Build stage
1. The Build Spec File Path is the relative location in your repo of the build_spec.yaml . Leave the default, for this example.
2. For the Primary Code Repository follow the below steps
   - Select connection type as GitHub
   - Select the external connection you created above
   - Give the HTTPS URL to the repo which contains your application.
   - Select main branch.

#### Create a Container Registry repository
Create a [Container Registry repository](https://docs.oracle.com/en-us/iaas/Content/Registry/Tasks/registrycreatingarepository.htm) for the `hello-world` container image built in the Managed Build stage.
1. You can name the repo: `java-docker-buildspec-sample-image`. So if you create the repository in the Ashburn region, the path is iad.ocir.io/TENANCY-NAMESPACE/java-docker-buildspec-sample-image
2. Set the repository access to private for security reasons. You can add policies to manage access for the same.
   
### Create a DevOps Artifact for your container image repository
Reference :  https://docs.oracle.com/en-us/iaas/Content/devops/using/containerimage_repository_artifact.htm

The version of the container image that will be delivered to the OCI repository is defined by a
[parameter](https://docs.oracle.com/en-us/iaas/Content/devops/using/configuring_parameters.htm) in the Artifact URI that
matches a Build Spec exported variable or Build Pipeline parameter name.

Create a DevOps Artifact to point to the Container Registry repository location you just created above. Enter the 
information for the Artifact location:
1. Name: `java-docker-buildspec-sample-artifact`
1. Type: Container image repository
1. Path: `REGION/TENANCY-OBJECT-STORAGE-NAMESPACE/java-docker-buildspec-sample-image`
1. Replace parameters: Yes

<img src="./assets/Create DevOps Artifact.png" />

Required policies must be added in the root compartment for the Container Registry repository and DevOps Artifact resource.
1. Provide access to OCIR to deliver artifacts : ```Allow dynamic-group dg-with-devops-resources to manage repos in tenancy```
2. Provide access to read deploy artifacts in deliver artifact stage : ```Allow dynamic-group dg-with-devops-resources to manage devops-family in tenancy```

### Add a Deliver Artifacts stage
Let's add a **Deliver Artifacts** stage to your Build Pipeline to deliver the `java-docker-buildspec-sample` container image to an OCI repository.

The Deliver Artifacts stage **maps** the ouput Artifacts from the Managed Build stage with the version to deliver to OCI Container Registry, through the DevOps Artifact Resource.

Add a **Deliver Artifacts** stage to your Build Pipeline after the **Managed Build** stage. To configure this stage:
1. In your Deliver Artifacts stage, choose `Select Artifact`
   <img src="./assets/Final Details in Deliver Artifact Stage.png" />
1. From the list of artifacts select the `java-docker-buildspec-sample-artifact` artifact that you created above
   <img src="./assets/Selecting DevOps Artifact resource.png" />
1. In the next section, you'll assign the  container image outputArtifact from the `build_spec.yaml` to the DevOps project artifact. For the "Build config/result Artifact name" enter: `hello_world_image`
   <img src="./assets/Final Details in Deliver Artifact Stage.png" />


### Run your Build in OCI DevOps

#### From your Build Pipeline, choose `Manual Run`
Use the Manual Run button to start a Build Run

Manual Run will use the Primary Code Repository, will start the Build Pipeline, first running the Managed Build stage, followed by the Deliver Artifacts stage.

After the Build Pipeline execution is complete, we can view the container image stored in the OCI Container Registry.


