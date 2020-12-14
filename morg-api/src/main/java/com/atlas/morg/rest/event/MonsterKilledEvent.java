package com.atlas.morg.rest.event;

public record MonsterKilledEvent(int worldId, int channelId, int mapId, int uniqueId, int monsterId, int x, int y,
                                 int killerId) {
}
