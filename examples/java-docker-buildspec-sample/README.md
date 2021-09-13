# Build Spec File Example using Java and Docker

This is an example Hello World project using Java SE JDK 11 and Docker. With the [OCI DevOps service](https://www.oracle.com/devops/devops-service/) and this project, you'll be able to build this application, build a Docker image and store it in [OCI Container Registry.](https://docs.oracle.com/en-us/iaas/Content/Registry/Concepts/registryoverview.htm)

In this example, you'll build a container image of this Java Hello World app, and store your built container in the OCI Container Registry,  all using the OCI DevOps service!

Let's go!

## Building the application locally

### Download the repo
The first step to get started is to download the repository to your local workspace

```shell
git clone git@github.com:AmodhShenoy/java-docker-buildspec-sample.git
cd java-docker-buildspec-sample
```

### Install and run the application

1. Install Java SE SDK 11 on your system: https://java.com/en/download/help/download_options.html
2. Compile the app: 
   
   ```javac src/com/sample/HelloWorld.java```
3. Run the app:
   
    ```java src/com/sample/HelloWorld```
4. To verify, make sure the string "Hello World" is printed in the shell.

### Build a container image for the app
You can locally build a container image using docker (or your favorite container image builder), to verify that you can run the app within a container.

```
docker build -t hello-world:1.0 .
```

Verify that your image was built, with `docker images`

Next run your local container and confirm you can access the app running in the container
```
docker run --rm -d --name hello-world:latest
```

The string "Hello World" must be printed in the shell if your container has been built successfully.

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
   - Finally, allow Build Service (dynamic group with DevOps Resources) to use a PAT secret by writing a policy in the root compartment as: ``` Allow dynamic-group dg-with-devops-resources to manage secret-family in tenancy```

### Setup your Build Pipeline
Create a new Build Pipeline to build, test and deliver artifacts from your GitHub Repository.

### Managed Build stage
In your Build Pipeline, first add a Managed Build stage
1. The Build Spec File Path is the relative location in your repo of the build_spec.yaml . Leave the default, for this example.
2. For the Primary Code Repository follow the below steps
    - Select connection type as GitHub
    - Select the external connection you created above
    - Give the URL to the repo which contains your application.
    - Select main branch.
    
### Create a Container Registry repository
Create a [Container Registry repository](https://docs.oracle.com/en-us/iaas/Content/Registry/Tasks/registrycreatingarepository.htm) for the `hello-world` container image built in the Managed Build stage.
1. You can name the repo: `java-docker-buildspec-sample-image`. So if you create the repository in the Ashburn region, the path is iad.ocir.io/TENANCY-NAMESPACE/java-docker-buildspec-sample-image
2. Set the repository access to public so that you can pull the container image without authorization from OKE. Under "Actions", choose `Change to public`.


### Create a DevOps Artifact for your container image repository
The version of the container image that will be delivered to the OCI repository is defined by a [parameter](https://docs.oracle.com/en-us/iaas/Content/devops/using/configuring_parameters.htm) in the Artifact URI that matches a Build Spec exported variable or Build Pipeline parameter name.

Create a DevOps Artifact to point to the Container Registry repository location you just created above. Enter the information for the Artifact location:
1. Name: `java-docker-buildspec-sample-artifact`
1. Type: Container image repository
1. Path: `iad.ocir.io/TENANCY-NAMESPACE/java-docker-buildspec-sample-image`
1. Replace parameters: Yes

### Add a Deliver Artifacts stage
Let's add a **Deliver Artifacts** stage to your Build Pipeline to deliver the `java-docker-buildspec-sample` container image to an OCI repository.

The Deliver Artifacts stage **maps** the ouput Artifacts from the Managed Build stage with the version to deliver to a DevOps Artifact resource, and then to the OCI repository.

Add a **Deliver Artifacts** stage to your Build Pipeline after the **Managed Build** stage. To configure this stage:
1. In your Deliver Artifacts stage, choose `Select Artifact`
1. From the list of artifacts select the `java-docker-buildspec-sample-artifact` artifact that you created above
1. In the next section, you'll assign the  container image outputArtifact from the `build_spec.yaml` to the DevOps project artifact. For the "Build config/result Artifact name" enter: `hello_world_image`


### Run your Build in OCI DevOps

#### From your Build Pipeline, choose `Manual Run`
Use the Manual Run button to start a Build Run

Manual Run will use the Primary Code Repository, will start the Build Pipeline, first running the Managed Build stage, followed by the Deliver Artifacts stage.

After the Build Pipeline execution is complete, we can view the container image stored in the OCI Container Registry, which can then be pulled to local workspace (Under ```Actions``` , choose ``` Copy Pull Command```).


