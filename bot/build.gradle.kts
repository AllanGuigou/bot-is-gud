import com.google.protobuf.gradle.generateProtoTasks
import com.google.protobuf.gradle.id
import com.google.protobuf.gradle.ofSourceSet
import com.google.protobuf.gradle.plugins
import com.google.protobuf.gradle.protobuf
import com.google.protobuf.gradle.protoc

val grpcVersion = "1.39.0"
val grpcKotlinVersion = "1.2.1"
val protobufVersion = "3.19.2"

plugins {
    kotlin("jvm") version "1.6.10"
    id("com.github.johnrengelman.shadow") version "2.0.2"
    id("com.google.protobuf") version "0.8.18"
}

group = "org.example"
version = "1.0-SNAPSHOT"

repositories {
    mavenLocal()
    google()
    jcenter()
    mavenCentral()
}

val ktlint by configurations.creating

dependencies {
    implementation(kotlin("stdlib-jdk8"))
    implementation("org.jetbrains.kotlin:kotlin-reflect:1.6.10")
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core")
    implementation("org.slf4j:slf4j-simple:1.7.30")
    implementation("dev.kord:kord-core:0.7.4")
    implementation("com.cronutils:cron-utils:9.1.0")

    implementation("io.grpc:grpc-kotlin-stub:$grpcKotlinVersion")
    implementation("io.grpc:grpc-protobuf:$grpcVersion")
    implementation("com.google.protobuf:protobuf-kotlin:$protobufVersion")
    runtimeOnly("io.grpc:grpc-netty-shaded:$grpcVersion")

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

sourceSets {
    main {
        proto {
            srcDir("build/generated/source/proto/main/java")
            srcDir("build/generated/source/proto/main/kotlin")
            srcDir("build/generated/source/proto/main/grpc")
            srcDir("build/generated/source/proto/main/grpckt")
        }
    }
}

protobuf {
    protoc {
        artifact = "com.google.protobuf:protoc:$protobufVersion"
    }
    plugins {
        id("grpc") {
            artifact = "io.grpc:protoc-gen-grpc-java:$grpcVersion"
        }
        id("grpckt") {
            artifact = "io.grpc:protoc-gen-grpc-kotlin:$grpcKotlinVersion:jdk7@jar"
        }
    }
    generateProtoTasks {
        all().forEach {
            it.plugins {
                id("grpc")
                id("grpckt")
            }
            it.builtins {
                id("kotlin")
            }
        }
    }
}

tasks {
    compileKotlin {
        kotlinOptions.jvmTarget = "1.8"
        kotlinOptions.freeCompilerArgs += "-opt-in=kotlin.RequiresOptIn"
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
