package com.atlas.morg.rest;

import javax.ws.rs.Consumes;
import javax.ws.rs.GET;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

import com.atlas.morg.rest.attribute.MonsterAttributes;
import com.atlas.morg.rest.processor.MonsterProcessor;

import rest.InputBody;

@Path("worlds")
public class WorldResource {
   @GET
   @Path("/{worldId}/channels/{channelId}/maps/{mapId}/monsters")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response getMonstersInMap(@PathParam("worldId") Integer worldId, @PathParam("channelId") Integer channelId,
                                      @PathParam("mapId") Integer mapId) {
      return MonsterProcessor.getMonstersInMap(worldId, channelId, mapId).build();
   }

   @POST
   @Path("/{worldId}/channels/{channelId}/maps/{mapId}/monsters")
   @Consumes(MediaType.APPLICATION_JSON)
   @Produces(MediaType.APPLICATION_JSON)
   public Response createMonsterInMap(@PathParam("worldId") Integer worldId, @PathParam("channelId") Integer channelId,
                                      @PathParam("mapId") Integer mapId,
                                      InputBody<MonsterAttributes> inputBody) {
      return MonsterProcessor.createMonster(worldId, channelId, mapId, inputBody.attributes()).build();
   }
}
