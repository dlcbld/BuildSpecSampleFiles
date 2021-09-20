# Build Spec File Example using Java (Gradle) and Creating a Jar file

This is an example Hello World project using Java SE JDK 11 and creating a jar file. With the [OCI DevOps service](https://www.oracle.com/devops/devops-service/) and this project, you'll be able to build this application,create a jar file and store it in [OCI Artifact Registry.](https://docs.oracle.com/en-us/iaas/artifacts/using/overview.htm)

In this example, We'll compile Java Hello World app by using gradle and  build a jar file, and store your built file in the OCI Artifact Registry,  all using the OCI DevOps service!

Let's go!

## Building the application locally

### Download the repo
The first step to get started is to download the repository to your local workspace

```shell
git clone https://github.com/anu-jha/java-gradle.git
cd java-gradle
```

### Install and run the application

1. Install Java SE SDK 11 on your system: https://java.com/en/download/help/download_options.html
2. Compile the app: 

   ```gradle compileJava```
3. Run the app:
   
    ```gradle run```
4. To verify, make sure the string "Hello World!" is printed in the shell.

### Build a Jar file 
1. We can locally build a jar file.

  ```
  gradle jar
  ```
 This command create a new folder with name "app" and jar file will be store in **app/build/libs**.
 
2. Verify that jar file was built, with name `app`, to run jar file we can simply run jar file.

```shell
cd app/build/libs
java -jar app.jar
```

3. To verify, make sure the string "Hello World!" is printed in the shell.


 ***Note :-***
  To remove pre build jar file we can use command 
  ```
   gradle clean
  ```
## Building the application in OCI DevOps
Now that We've seen we can locally build this app, let's try this out through OCI DevOps Build service.

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
    
### Create a Artifact Registry repository
Create a [Artifact Registry repository](https://docs.oracle.com/en-us/iaas/artifacts/using/manage-repos.htm#create-repo) for the `build-jar` file built in the Managed Build stage.

### Create a DevOps Artifact for your Artifact repository
The version of the jar file that will be delivered to the OCI repository is defined by a [parameter](https://docs.oracle.com/en-us/iaas/Content/devops/using/configuring_parameters.htm) in the Artifact URI that matches a Build Spec exported variable or Build Pipeline parameter name.

Create a DevOps Artifact to point to the Artifact Registry repository location you just created above. Enter the information for the Artifact location:
1. Name: `java-gradle-buildspec-sample-artifact`
2. Type: General Artifact
3. Artifact source: Artifact Registry repository
4. Version: ${version} (assign some parameter eg:- 1.0, 2.0)
### Add a Deliver Artifacts stage
Let's add a **Deliver Artifacts** stage to your Build Pipeline to deliver the `build_jar` jar file to an OCI repository.

The Deliver Artifacts stage **maps** the ouput Artifacts from the Managed Build stage with the version to deliver to a DevOps Artifact resource, and then to the OCI repository.

Add a **Deliver Artifacts** stage to your Build Pipeline after the **Managed Build** stage. To configure this stage:
1. In your Deliver Artifacts stage, choose `Select Artifact`
1. From the list of artifacts select the `java-jar-buildspec-sample-artifact` artifact that you created above
1. In the next section, We'll assign the  jar file outputArtifact from the `build_spec.yaml` to the DevOps project artifact. For the "Build config/result Artifact name" enter: `build_jar`


### Run your Build in OCI DevOps

#### From your Build Pipeline, choose `Manual Run`
Use the Manual Run button to start a Build Run

Manual Run will use the Primary Code Repository, will start the Build Pipeline, first running the Managed Build stage, followed by the Deliver Artifacts stage.

After the Build Pipeline execution is complete, we can view the jar file stored in the OCI Artifact Registry, which can then be download to local workspace. 

