package com.example.controller;

import com.example.model.User;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.tags.Tag;
import org.springframework.web.bind.annotation.*;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicLong;

@RestController
@RequestMapping("/api/users")
@Tag(name = "用户管理", description = "用户的增删改查接口")
public class UserController {

    private final Map<Long, User> usersDb = new ConcurrentHashMap<>();
    private final AtomicLong idGenerator = new AtomicLong(1);

    public UserController() {
        // 初始化示例数据
        User user1 = new User();
        user1.setId(idGenerator.getAndIncrement());
        user1.setName("张三");
        user1.setEmail("zhangsan@example.com");
        user1.setRole("admin");
        usersDb.put(user1.getId(), user1);

        User user2 = new User();
        user2.setId(idGenerator.getAndIncrement());
        user2.setName("李四");
        user2.setEmail("lisi@example.com");
        user2.setRole("user");
        usersDb.put(user2.getId(), user2);
    }

    @GetMapping
    @Operation(summary = "获取用户列表", description = "返回系统中所有用户的信息")
    public List<User> listUsers() {
        return new ArrayList<>(usersDb.values());
    }

    @GetMapping("/{id}")
    @Operation(summary = "根据ID获取用户", description = "根据用户ID返回单个用户信息")
    public User getUser(
            @Parameter(description = "用户ID", required = true, example = "1")
            @PathVariable Long id) {
        return usersDb.get(id);
    }

    @PostMapping
    @Operation(summary = "创建用户", description = "创建一个新的用户")
    public User createUser(@RequestBody User user) {
        user.setId(idGenerator.getAndIncrement());
        usersDb.put(user.getId(), user);
        return user;
    }

    @PutMapping("/{id}")
    @Operation(summary = "更新用户", description = "根据ID更新用户信息")
    public User updateUser(
            @Parameter(description = "用户ID", required = true)
            @PathVariable Long id,
            @RequestBody User user) {
        user.setId(id);
        usersDb.put(id, user);
        return user;
    }

    @DeleteMapping("/{id}")
    @Operation(summary = "删除用户", description = "根据ID删除用户")
    public boolean deleteUser(
            @Parameter(description = "用户ID", required = true)
            @PathVariable Long id) {
        return usersDb.remove(id) != null;
    }
}
