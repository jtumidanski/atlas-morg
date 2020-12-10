package com.atlas.morg.processor;

import java.util.Collections;
import java.util.List;
import java.util.stream.Collectors;

import com.atlas.mrg.rest.attribute.MapCharacterAttributes;
import com.atlas.shared.rest.RestService;
import com.atlas.shared.rest.UriBuilder;

import rest.DataBody;
import rest.DataContainer;

public class MapProcessor {
   public static List<Integer> getCharacterIdsInMap(int worldId, int channelId, int mapId) {
      return UriBuilder.service(RestService.MAP_REGISTRY)
            .pathParam("worlds", worldId)
            .pathParam("channels", channelId)
            .pathParam("maps", mapId)
            .path("characters")
            .getRestClient(MapCharacterAttributes.class)
            .getWithResponse()
            .result()
            .map(DataContainer::getDataAsList)
            .orElse(Collections.emptyList())
            .stream()
            .map(DataBody::getId)
            .map(Integer::parseInt)
            .collect(Collectors.toList());
   }
}
