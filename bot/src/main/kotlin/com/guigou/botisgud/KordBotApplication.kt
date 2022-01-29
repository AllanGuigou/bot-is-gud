package com.guigou.botisgud

import com.guigou.botisgud.commands.Presence
import com.guigou.botisgud.commands.Typing
import com.guigou.botisgud.commands.register
import com.guigou.botisgud.services.grpc.PresenceServiceClient
import dev.kord.core.Kord
import io.grpc.ManagedChannelBuilder
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val token = System.getenv("DISCORD_TOKEN")
    val client = Kord(token)

    val builder = ManagedChannelBuilder.forTarget("localhost:50051").usePlaintext()
    // TODO: research the usage of an executor for the builder
    val presenceServiceClient = PresenceServiceClient(builder.build())

    // https://elizarov.medium.com/coroutine-context-and-scope-c8b255d59055
    client.register(Typing(), this) // TODO: determine how to initialize class within extension function
    client.register(Presence(presenceServiceClient), this)

    client.login()
}
