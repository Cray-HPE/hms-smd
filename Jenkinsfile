@Library('dst-shared@release/shasta-1.4') _

dockerBuildPipeline {
        repository = "cray"
        imagePrefix = "cray"
        app = "smd"
        name = "hms-smd"
        description = "Cray Hardware State Manager Service"
        dockerfile = "Dockerfile.smd"
        slackNotification = ["", "", false, false, true, true]
        product = "csm"
}
