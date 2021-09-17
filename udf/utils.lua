function remove(rec)
    info("##############starting remove UDF on ID "  .. tostring(rec.ID) .. " with Timestamp " .. tostring(rec.Timestamp))
    aerospike:remove(rec)
    info("##############done remove UDF")
end