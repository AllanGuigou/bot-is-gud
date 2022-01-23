package com.guigou.botisgud.services.nickname

import org.jetbrains.exposed.dao.IntEntity
import org.jetbrains.exposed.dao.IntEntityClass
import org.jetbrains.exposed.dao.id.EntityID

class NicknameEntity(id: EntityID<Int>) : IntEntity(id) {
    companion object: IntEntityClass<NicknameEntity>(Nicknames)

    var userId by Nicknames.userId
    var guildId by Nicknames.guildId
}
