plugins {
    kotlin("jvm") version "1.4.0"
    id("com.github.johnrengelman.shadow") version "2.0.2"
}

group = "org.example"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
    jcenter()
    maven(url = "https://dl.bintray.com/kordlib/Kord")
}

val ktlint by configurations.creating

dependencies {
    implementation(kotlin("stdlib-jdk8"))
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core")
    implementation("org.slf4j:slf4j-simple:1.7.30")
    implementation("com.gitlab.kordlib.kord:kord-core:0.6.10")
    implementation("com.cronutils:cron-utils:9.1.0")

    implementation("org.jetbrains.exposed:exposed-core:0.25.1")
    implementation("org.jetbrains.exposed:exposed-dao:0.25.1")
    implementation("org.jetbrains.exposed:exposed-jdbc:0.25.1")
    implementation("org.jetbrains.exposed:exposed-java-time:0.25.1")
    implementation("org.postgresql:postgresql:42.2.2")

    ktlint("com.pinterest:ktlint:0.39.0")

    testImplementation(platform("org.junit:junit-bom:5.7.0"))
    testImplementation("org.junit.jupiter:junit-jupiter")
    testImplementation("io.mockk:mockk:1.10.0")
    testImplementation("com.willowtreeapps.assertk:assertk-jvm:0.22")
    testImplementation("org.xerial:sqlite-jdbc:3.30.1")
}

tasks {
    compileKotlin {
        kotlinOptions.jvmTarget = "1.8"
    }
    compileTestKotlin {
        kotlinOptions.jvmTarget = "1.8"
    }
    test {
        useJUnitPlatform()
    }
    withType<Jar> {
        manifest {
            attributes(mapOf("Main-Class" to "com.guigou.botisgud.KordBotApplicationKt"))
        }
    }
}

val outputDir = "${project.buildDir}/reports/ktlint/"
val inputFiles = project.fileTree(mapOf("dir" to "src", "include" to "**/*.kt"))

val ktlintCheck by tasks.creating(JavaExec::class) {
    inputs.files(inputFiles)
    outputs.dir(outputDir)

    description = "Check Kotlin code style."
    classpath = ktlint
    main = "com.pinterest.ktlint.Main"
    args = listOf("src/**/*.kt")
}

val ktlintFormat by tasks.creating(JavaExec::class) {
    inputs.files(inputFiles)
    outputs.dir(outputDir)

    description = "Fix Kotlin code style deviations."
    classpath = ktlint
    main = "com.pinterest.ktlint.Main"
    args = listOf("-F", "src/**/*.kt")
}
