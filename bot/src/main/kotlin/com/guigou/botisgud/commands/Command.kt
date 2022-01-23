package com.guigou.botisgud.commands

import dev.kord.core.Kord
import kotlinx.coroutines.CoroutineScope

interface Command {
    suspend fun register(client: Kord, scope: CoroutineScope)
}

suspend fun <T : Command> Kord.register(command: T, scope: CoroutineScope) {
    command.register(this, scope)
}
