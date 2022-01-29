package com.guigou.botisgud

import com.guigou.botisgud.commands.Activity
import com.guigou.botisgud.commands.Typing
import com.guigou.botisgud.commands.register
import dev.kord.core.Kord
import kotlinx.coroutines.runBlocking

fun main() = runBlocking {
    val token = System.getenv("DISCORD_TOKEN")
    val client = Kord(token)

    // https://elizarov.medium.com/coroutine-context-and-scope-c8b255d59055
    client.register(Typing(), this) // TODO: determine how to initialize class within extension function
    client.register(Activity(), this)

    client.login()
}
