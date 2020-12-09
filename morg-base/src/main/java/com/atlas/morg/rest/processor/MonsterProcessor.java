package com.atlas.morg.rest.processor;

import javax.ws.rs.core.Response;

import com.app.rest.util.stream.Collectors;
import com.app.rest.util.stream.Mappers;
import com.atlas.mis.attribute.MonsterDataAttributes;
import com.atlas.morg.MonsterRegistry;
import com.atlas.morg.rest.ResultObjectFactory;
import com.atlas.morg.rest.attribute.MonsterAttributes;
import com.atlas.shared.rest.RestService;
import com.atlas.shared.rest.UriBuilder;

import builder.ResultBuilder;
import rest.DataBody;
import rest.DataContainer;

public final class MonsterProcessor {
   private MonsterProcessor() {
   }

   public static ResultBuilder createMonster(int worldId, int channelId, int mapId, MonsterAttributes attributes) {
      return UriBuilder.service(RestService.MAP_INFORMATION)
            .pathParam("monsters", attributes.monsterId())
            .getRestClient(MonsterDataAttributes.class)
            .getWithResponse()
            .result()
            .map(DataContainer::getData)
            .map(DataBody::getAttributes)
            .map(data -> com.atlas.morg.processor.MonsterProcessor.createMonster(worldId, channelId, mapId,
                  attributes.monsterId(), attributes.x(), attributes.y(), attributes.fh(), attributes.stance(), attributes.team())
            )
            .map(ResultObjectFactory::create)
            .map(Mappers::singleCreatedResult)
            .orElse(new ResultBuilder(Response.Status.FORBIDDEN));
   }

   public static ResultBuilder getMonster(int uniqueId) {
      return MonsterRegistry.getInstance().getMonster(uniqueId)
            .map(ResultObjectFactory::create)
            .map(Mappers::singleCreatedResult)
            .orElse(new ResultBuilder(Response.Status.FORBIDDEN));
   }

   public static ResultBuilder getMonstersInMap(int worldId, int channelId, int mapId) {
      return MonsterRegistry.getInstance().getMonstersInMap(worldId, channelId, mapId).stream()
            .map(ResultObjectFactory::create)
            .collect(Collectors.toResultBuilder());
   }
}
