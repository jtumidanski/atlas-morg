package com.atlas.morg.rest.event;

public record MonsterEvent(int worldId, int channelId, int mapId, int uniqueId, int monsterId, Integer actorId,
                           MonsterEventType type) {
}
