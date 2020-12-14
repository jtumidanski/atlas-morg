package com.atlas.morg.rest.attribute;

import rest.AttributeResult;

public record MonsterAttributes(Integer worldId, Integer channelId, Integer mapId, Integer monsterId, Integer controlCharacterId,
                                Integer x, Integer y, Integer fh, Integer stance, Integer team, Integer hp)
      implements AttributeResult {
}
