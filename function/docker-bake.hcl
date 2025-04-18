group "default" {
    targets = ["main", "appservice"]
}

target "common" {
    context = "."
        dockerfile = "Dockerfile"
    args = {
        ALPINE_VERSION = "3.20"
    }
}

target "main" {
    inherits = ["common"]
    target = "final"
    tags = ["clowa/az-func-daily-quote:dev"]
}

target "appservice" {
    inherits = ["common"]
    target = "appservice"
    tags = ["clowa/az-func-daily-quote:appservice"]
}
