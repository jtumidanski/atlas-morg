package com.atlas.morg.processor;

import java.util.Collections;
import java.util.stream.Stream;

import com.atlas.mrg.constant.RestConstants;
import com.atlas.mrg.rest.attribute.MapCharacterAttributes;
import com.atlas.shared.rest.UriBuilder;

import rest.DataBody;
import rest.DataContainer;

public class MapProcessor {
   public static Stream<Integer> getCharacterIdsInMap(int worldId, int channelId, int mapId) {
      return UriBuilder.service(RestConstants.SERVICE)
            .pathParam("worlds", worldId)
            .pathParam("channels", channelId)
            .pathParam("maps", mapId)
            .path("characters")
            .getRestClient(MapCharacterAttributes.class)
            .getWithResponse()
            .result()
            .map(DataContainer::dataList)
            .orElse(Collections.emptyList())
            .stream()
            .map(DataBody::getId)
            .map(Integer::parseInt);
   }
}
