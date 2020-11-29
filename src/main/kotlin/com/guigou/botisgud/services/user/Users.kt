package com.guigou.botisgud.services.user

import org.jetbrains.exposed.dao.id.IntIdTable

object Users : IntIdTable() {
    val userId = long("userId")
    val zone = varchar("zone", 50)
}
