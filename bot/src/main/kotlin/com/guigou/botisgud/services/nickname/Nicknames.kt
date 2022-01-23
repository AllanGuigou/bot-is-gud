package com.guigou.botisgud.services.nickname

import org.jetbrains.exposed.dao.id.IntIdTable

object Nicknames : IntIdTable() {
    val userId = long("userId")
    val guildId = long("guildId")
}
