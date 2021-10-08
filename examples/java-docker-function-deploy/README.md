# Build Spec File Example using Java and Docker to deploy to a Function

This is an example Hello World project using Java SE JDK 11 and Docker. With the [OCI DevOps service](https://www.oracle.com/devops/devops-service/) 
and this project, you'll be able to build this application, build a Docker image and store it in 
[OCI Container Registry.](https://docs.oracle.com/en-us/iaas/Content/Registry/Concepts/registryoverview.htm) We will then trigger a deployment pipeline
which deploys this image to a function.

In this example, you'll build a container image of this Java Hello World app, store your built container in the OCI Container Registry, and trigger a 
function deployment all using the OCI DevOps service!

Let's go!

## Building the application locally

### Download the repo
The first step to get started is to download the repository to your local workspace

```shell
git clone git@github.com:dlcbld/BuildSpecSampleFiles.git
cd examples/java-docker-buildspec-sample
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

### Setting up your Function

#### Create a Container Registry repository
Create a [Container Registry repository](https://docs.oracle.com/en-us/iaas/Content/Registry/Tasks/registrycreatingarepository.htm) for the `hello-world` container image built in the Managed Build stage.
1. You can name the repo: `java-docker-buildspec-sample-image`. So if you create the repository in the Ashburn region, the path is iad.ocir.io/TENANCY-NAMESPACE/java-docker-buildspec-sample-image
2. Set the repository access to private for security reasons. You can add policies to manage access for the same.

#### Pushing your image to Container Registry repository
Reference: https://www.oracle.com/webfolder/technetwork/tutorials/obe/oci/registry/index.html#TagtheImageforPushing

a. Docker Login
   1. In the top-right corner of the Console, open the User menu (User menu), and then click User Settings.
   2. On the Auth Tokens page, click Generate Token. Enter the details, and click on create.
   3. Copy the auth token immediately to a secure location from where you can retrieve it later, because you won't see the auth token again in the Console.
   4. In a terminal window on the client machine running Docker, log in to Oracle Cloud Infrastructure Registry by entering:
      
      ```
      docker login <region-key>.ocir.io
      ```
      where <region-key> is the key for the Oracle Cloud Infrastructure Registry region you're using. For example, ```iad```. See the Availability by Region topic in the Oracle Cloud Infrastructure Registry documentation. 
   5. When prompted, enter your username in the format <tenancy-namespace>/<username>. For example, ansh81vru1zp/jdoe@acme.com. If your tenancy is federated with Oracle Identity Cloud Service, use the format <tenancy-namespace>/oracleidentitycloudservice/<username>.
   6. When prompted, enter the auth token you copied earlier as the password.

b. Pushing image through CLI
   1. In a terminal, give a tag to the image that you're going to push to OCI Registry by entering:
      
      ```
      docker tag hello-world:latest <region-key>.ocir.io/<tenancy-namespace>/<repo-name>:<tag>
      ```
      where:
      - *<region-key>* is the key for the Oracle Cloud Infrastructure Registry region you're using. For example, ```iad```. See the Availability by Region topic in the Oracle Cloud Infrastructure Registry documentation.
      ocir.io is the Oracle Cloud Infrastructure Registry name.
      - *<tenancy-namespace>* is the auto-generated Object Storage namespace string of the tenancy (as shown on the Tenancy Information page) to which you want to push the image. For example, the namespace of the acme-dev tenancy might be ansh81vru1zp. Note that your user must have access to the tenancy.
      <repo-name> is the name of the target repository to which you want to push the image (for example, helloworld). Note that you'll usually specify a repository that already exists.
      - *< tag>* is an image tag you want to give the image in Oracle Cloud Infrastructure Registry (for example, latest).
      
      For example:
      ```
      docker tag hello-world:latest iad.ocir.io/ansh81vru1zp/java-docker-buildspec-sample-image:latest
      ```
   2. In a terminal, push the Docker image from the client machine to Oracle Cloud Infrastructure Registry by entering:
      
      ```
      docker push <region-key>.ocir.io/<tenancy-namespace>/<repo-name>:<tag>
      ```


#### Setting up a VCN
1. In the menu, go to Networking and select Virtual Cloud Networks.
2. Click on Create VCN and fill in the necessary details as below:
   
<img src="./assets/Create VCN.png" />

3. Open up your created VCN and click on Create Subnet. Fill in the necessary details
as below and create your private subnet.
   
<img src="./assets/Create Subnet.png" />

#### Creating an Application
1. In the menu, under Featured you will find Applications. Go to the page and select Create Application
2. Fill in the necessary details and select the VCN and subnet created in the earlier step. Click on Create.

<img src="./assets/Create Application.png" />

#### Creating the Function 
Move into your created Application and click on Create Function. Fill in the necessary details
and select your container image repository created in the earlier step and the image uploaded inside it.
   
With this you have successfully created a Function on OCI DevOps. Now let us try to automate this
Function deploy/update part through OCI Build Service.

### Setting up an Environment
1. Create a [DevOps Project](https://docs.oracle.com/en-us/iaas/Content/devops/using/devops_projects.htm) or use and an existing project.
2. In your DevOps project, go to the Environments section and select Create Environment.
3. Fill in the necessary details and select ```Create Environment for a Function```.
4. In the next step, choose your Application and the Function created in the previous steps.

<img src="./assets/Create Environment.png" />

### Create External Connection to your Git repository

1. In your DevOps project, create an External Connection to your GitHub repository which holds your application.
   - Create a Personal Access Token (PAT): https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token
   - In the OCI Console, Go to Identity & Security -> Vault and create a [Vault]( https://docs.oracle.com/en-us/iaas/Content/KeyManagement/Concepts/keyoverview.htm) in compartment of your own choice.
   - Create a Master Key that will be used to encrypt the PATs.
   - Select Secrets from under Resources and create a secret using PAT obtained from GitHub account.
   - Make a note of the OCID of the secret.
   - Now, go to the desired project and select External Connection from the resources.
   - Select type as GitHub and provide OCID of the secret under Personal Access Token.
   - Finally, allow BuildPipeline (dynamic group with DevOps Resources) to use a PAT secret by writing a policy in the root compartment as: ``` Allow dynamic-group dg-with-devops-resources to manage secret-family in tenancy```
   
### Create a DevOps Artifact for your container image repository
Reference :  https://docs.oracle.com/en-us/iaas/Content/devops/using/containerimage_repository_artifact.htm

The version of the container image that will be delivered to the OCI repository is defined by a [parameter](https://docs.oracle.com/en-us/iaas/Content/devops/using/configuring_parameters.htm) in the Artifact URI that matches a Build Spec exported variable or Build Pipeline parameter name.

Create a DevOps Artifact to point to the Container Registry repository location you just created above. Enter the information for the Artifact location:
1. Name: `java-docker-buildspec-sample-artifact`
1. Type: Container image repository
1. Path: `REGION/TENANCY-OBJECT-STORAGE-NAMESPACE/java-docker-buildspec-sample-image`
1. Replace parameters: Yes

<img src="./assets/Create DevOps Artifact.png" />

Make sure you have added required policies for both your DevOps Artifact and your Container Repository according to [this](https://docs.oracle.com/en-us/iaas/Content/Identity/Reference/registrypolicyreference.htm)

### Setting up your Deployment Pipeline
Create a new Deployment Pipeline to deploy the Function that we just created.

### Function Deployment Stage
1. In your Deployment Pipeline, select the default Functions stage under Deploy. 
2. Enter the necessary details and select the Environment, Function to run and the artifact created as well.

### Setting up your Build Pipeline
Create a new Build Pipeline to build, test and deliver artifacts from your GitHub Repository.

### Managed Build stage
In your Build Pipeline, first add a Managed Build stage
1. The Build Spec File Path is the relative location in your repo of the build_spec.yaml . Leave the default, for this example.
2. For the Primary Code Repository follow the below steps
   - Select connection type as GitHub
   - Select the external connection you created above
   - Give the HTTPS URL to the repo which contains your application.
   - Select main branch.

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
   
### Trigger Deployment Stage
1. In your BuildPipeline, again go to add stage and select Trigger Deployment Stage.
2. Enter the necessary details and select the Deployment Pipeline that we just created which maps to our target Function.

### Putting it all together

From your Build Pipeline, choose `Manual Run`.
Manual Run will use the Primary Code Repository, will start the Build Pipeline, first running the Managed Build stage, followed by the Deliver Artifacts stage.At the
end of this stage, the image is newly built and is updated in the container repository we created. Then it triggers a deployment pipeline which then deploys/updates
the function with the new image from our container repository.

